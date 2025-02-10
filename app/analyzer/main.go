package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/sptGabriel/investment-analyzer/app/analyzer/api"
	v1 "github.com/sptGabriel/investment-analyzer/app/analyzer/api/v1"
	"github.com/sptGabriel/investment-analyzer/domain/challenge"
	"github.com/sptGabriel/investment-analyzer/domain/reports"
	"github.com/sptGabriel/investment-analyzer/extensions/gbdb"
	"github.com/sptGabriel/investment-analyzer/extensions/gblib"
	"github.com/sptGabriel/investment-analyzer/gateways/assets"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres"
	assetsRepo "github.com/sptGabriel/investment-analyzer/gateways/postgres/assets"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/auditlogs"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/portfolios"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/prices"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/settings"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/trades"
	"github.com/sptGabriel/investment-analyzer/interceptors"
	"github.com/sptGabriel/investment-analyzer/telemetry"
	"github.com/sptGabriel/investment-analyzer/telemetry/logging"
)

var (
	AppName     = "investment-analyzer"
	BuildCommit = "undefined"
	BuildTag    = "undefined"
	BuildTime   = "undefined"
)

type envs struct {
	AppName string `conf:"env:APP_NAME,required"`

	GracefulShutdown time.Duration `conf:"env:GRACEFUL_SHUTDOWN_TIMEOUT,default:20s"`

	ServerAddr         string        `conf:"env:SERVER_ADDR,default:0.0.0.0:3000"`
	ServerReadTimeout  time.Duration `conf:"env:SERVER_READ_TIMEOUT,default:30s"`
	ServerWriteTimeout time.Duration `conf:"env:SERVER_WRITE_TIMEOUT,default:30s"`

	// Postgres
	DBName                  string `conf:"env:DATABASE_NAME,default:investment_analyzer"`
	DBUser                  string `conf:"env:DATABASE_USER,default:postgres"`
	DBPassword              string `conf:"env:DATABASE_PASSWORD,default:postgres"`
	DBHost                  string `conf:"env:DATABASE_HOST_DIRECT,default:localhost"`
	DBPort                  string `conf:"env:DATABASE_PORT_DIRECT,default:5433"`
	DBPoolMinSize           int32  `conf:"env:DATABASE_POOL_MIN_SIZE,default:2"`
	DBPoolMaxSize           int32  `conf:"env:DATABASE_POOL_MAX_SIZE,default:10"`
	DBPoolMaxConnLifetime   string `conf:"env:DATABASE_POOL_MAX_CONN_LIFETIME"`
	DBPoolMaxConnIdleTime   string `conf:"env:DATABASE_POOL_MAX_CONN_IDLE_TIME"`
	DBPoolHealthCheckPeriod string `conf:"env:DATABASE_POOL_HEALTHCHECK_PERIOD"`
	DBSSLMode               string `conf:"env:DATABASE_SSLMODE,default:disable"`
	DBSSLRootCert           string `conf:"env:DATABASE_SSL_ROOTCERT"`
	DBSSLCert               string `conf:"env:DATABASE_SSL_CERT"`
	DBSSLKey                string `conf:"env:DATABASE_SSL_KEY"`
}

func (e envs) PGAddress() string {
	if e.DBSSLMode == "" {
		e.DBSSLMode = "disable"
	}

	address := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		e.DBUser, e.DBPassword, e.DBHost, e.DBPort, e.DBName, e.DBSSLMode)

	return address
}

func main() {
	logger := logging.New(os.Stderr).With(
		zap.String("build_commit", BuildCommit),
		zap.String("build_time", BuildTime),
		zap.String("build_tag", BuildTag),
		zap.Int("go_max_procs", runtime.GOMAXPROCS(0)),
		zap.Int("runtime_num_cpu", runtime.NumCPU()),
	)

	var cfg envs
	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			logger.Info(help)
			return
		}

		logger.Panic("error parsing the configuration", zap.Error(err))
	}

	logger.Info(help)

	if err = startApp(logger, cfg); err != nil {
		logger.Panic("running analyzer API", zap.Error(err))
	}
}

func startApp(logger *zap.Logger, cfg envs) error {
	ctx := logging.WithContext(context.Background(), logger)

	dbConn, err := postgres.New(
		cfg.PGAddress(),
		cfg.DBPoolMinSize,
		cfg.DBPoolMaxSize,
	)
	if err != nil {
		return fmt.Errorf("setuping postgres: %w", err)
	}

	db := gbdb.NewDatabase(dbConn)
	tx := gbdb.NewTransactioner(db)

	tradeRepository := trades.New(db)
	pricesRepository := prices.New(db)
	assetsRepository := assetsRepo.New(db)
	portfoliosRepository := portfolios.New(db)
	settingsRepository := settings.New(db)
	auditRepository := auditlogs.New(db)

	csvImported := false
	if err := db.With(ctx, func(ctx context.Context) error {
		var innerErr error
		csvImported, innerErr = settingsRepository.IsCsvImported(ctx)
		return innerErr
	}); err != nil {
		return fmt.Errorf("on csv imported: %w", err)
	}

	if !csvImported {
		csvReader, err := assets.New()
		if err != nil {
			return fmt.Errorf("on csv reader setup: %w", err)
		}

		setupChallengeUC := challenge.NewSetupChallenge(
			db, tx, csvReader, settingsRepository, assetsRepository,
			pricesRepository, portfoliosRepository, tradeRepository)

		if _, err = setupChallengeUC.Execute(ctx, challenge.SetupChallengeInput{}); err != nil {
			return fmt.Errorf("on import csv data: %w", err)
		}
	}

	router, err := api.NewServer(logger)
	if err != nil {
		return fmt.Errorf("bulding api server: %w", err)
	}

	reportGeneratorUC := reports.NewGenerateReportUC(db, tradeRepository, pricesRepository, portfoliosRepository)

	reportGeneratorUCWithAudit := gblib.New(
		reportGeneratorUC,
		gblib.WithDB(db),
		interceptors.AuditInterceptor(auditRepository),
	)

	apiV1 := v1.API{
		ReportHandler: v1.NewReportHandler(reportGeneratorUCWithAudit),
	}

	apiV1.Routes(router)

	apiServer := http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}

	metricsServer, err := telemetry.NewMetricsServer()
	if err != nil {
		return fmt.Errorf("building metrics server: %w", err)
	}

	// Graceful Shutdown
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	g, gCtx := errgroup.WithContext(signalCtx)

	g.Go(func() error {
		logger.Info("starting api server")

		logger.Info("api started",
			zap.String("address", apiServer.Addr),
			zap.Duration("read_timeout", apiServer.ReadTimeout),
			zap.Duration("write_timeout", apiServer.WriteTimeout),
		)

		return apiServer.ListenAndServe()
	})

	g.Go(func() error {
		logger.Info("metrics server started",
			zap.String("address", metricsServer.Addr),
			zap.Duration("read_timeout", metricsServer.ReadTimeout),
			zap.Duration("write_timeout", metricsServer.WriteTimeout),
		)

		return metricsServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		logger.Info("interrupt signal received")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
		defer cancel()

		var errs error
		logger.Info("closing http server")
		if err := apiServer.Shutdown(timeoutCtx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to stop api server: %w", err))
		}

		logger.Info("closing metrics server")
		if err := metricsServer.Shutdown(timeoutCtx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to stop metrics server: %w", err))
		}

		return errs
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("on finishing app", zap.Error(err))
	}

	stop()

	logger.Info("bye")

	return nil
}

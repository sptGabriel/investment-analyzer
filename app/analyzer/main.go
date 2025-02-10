package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"

	"github.com/sptGabriel/investment-analyzer/domain/challenge"
	"github.com/sptGabriel/investment-analyzer/domain/reports"
	"github.com/sptGabriel/investment-analyzer/extensions/gbdb"
	"github.com/sptGabriel/investment-analyzer/gateways/assets"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres"
	assetsRepo "github.com/sptGabriel/investment-analyzer/gateways/postgres/assets"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/portfolios"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/prices"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/settings"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/trades"
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
	ctx := context.Background()

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

	generateReportUC := reports.NewGenerateReportUC(
		db, tradeRepository, pricesRepository, portfoliosRepository,
	)

	start, _ := time.Parse("2006-01-02 15:04:05", "2021-03-01 10:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2021-03-07 17:50:00")
	interval := 10 * time.Minute

	output, err := generateReportUC.Execute(ctx, reports.GenerateReportInput{
		PortfolioID: "408186c6-b76a-4ad6-8d4a-9ace3762b997",
		StartDate:   start,
		EndDate:     end,
		Interval:    interval,
	})

	fmt.Println(err, output)

	return nil
}

package api

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/sptGabriel/investment-analyzer/extensions/gbhttp/middlewares"
	"github.com/sptGabriel/investment-analyzer/telemetry"
)

var defaultTimeout = 15 * time.Second

func NewServer(logger *zap.Logger) (*chi.Mux, error) {
	redMetrics, err := telemetry.NewRedMetricsMiddleware()
	if err != nil {
		return nil, fmt.Errorf("creating red metrics middleware: %w", err)
	}

	router := chi.NewRouter()
	router.Use(
		middlewares.WithLogger(logger),
		middlewares.RequestDataMiddleware,
		redMetrics.Handle(),
		chimid.Recoverer,
		chimid.Timeout(defaultTimeout),
	)

	return router, nil
}

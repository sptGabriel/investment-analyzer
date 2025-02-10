package api

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/sptGabriel/investment-analyzer/extensions/http/middlewares"
	"github.com/sptGabriel/investment-analyzer/telemetry"
)

func NewServer(logger *zap.Logger) (*chi.Mux, error) {
	redMetrics, err := telemetry.NewRedMetricsMiddleware()
	if err != nil {
		return nil, fmt.Errorf("creating red metrics middleware: %w", err)
	}

	router := chi.NewRouter()
	router.Use(
		middlewares.WithLogger(logger),
		redMetrics.Handle(),
		chimid.Recoverer,
	)

	return router, nil
}

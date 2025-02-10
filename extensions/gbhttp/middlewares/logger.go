package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/sptGabriel/investment-analyzer/telemetry/logging"
)

func WithLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(logging.WithContext(r.Context(), logger))
			next.ServeHTTP(w, r)
		})
	}
}

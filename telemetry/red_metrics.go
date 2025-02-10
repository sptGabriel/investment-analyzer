package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/sptGabriel/investment-analyzer/extensions/utils"
)

type RedMetricsMiddleware struct {
	httpServerRequestDuration metric.Float64Histogram
}

func (m *RedMetricsMiddleware) Handle() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			dw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(dw, r)

			ctx := r.Context()
			attrs := []attribute.KeyValue{
				semconv.HTTPResponseStatusCode(dw.Status()),
				semconv.HTTPRoute(chi.RouteContext(ctx).RoutePattern()),
				semconv.HTTPRequestMethodKey.String(r.Method),
				semconv.NetworkProtocolName(r.Proto),
			}

			if r.URL != nil {
				attrs = append(attrs, semconv.URLScheme(r.URL.Scheme))
			}

			m.httpServerRequestDuration.Record(ctx, utils.TimeSince(start), metric.WithAttributeSet(attribute.NewSet(attrs...)))
		}

		return http.HandlerFunc(fn)
	}
}

func NewRedMetricsMiddleware() (*RedMetricsMiddleware, error) {
	meter := otel.GetMeterProvider().Meter("investment-analyzer")

	hist, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("The duration of the HTTP request."),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
		metric.WithUnit("s"),
	)

	if err != nil {
		return nil, fmt.Errorf("creating float64 histogram: %w", err)
	}

	return &RedMetricsMiddleware{httpServerRequestDuration: hist}, nil
}

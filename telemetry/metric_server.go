package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsConfig struct {
	Addr         string        `conf:"env:METRICS_SERVER_ADDR,default:0.0.0.0:3001"`
	ReadTimeout  time.Duration `conf:"env:METRICS_SERVER_READ_TIMEOUT,default:30s"`
	WriteTimeout time.Duration `conf:"env:METRICS_SERVER_WRITE_TIMOUT,default:30s"`
}

func NewMetricsServer() (*http.Server, error) {
	var cfg MetricsConfig

	_, err := conf.Parse("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing metrics server config: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Handler:      mux,
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}, nil
}

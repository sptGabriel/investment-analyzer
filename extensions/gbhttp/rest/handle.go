package rest

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/sptGabriel/investment-analyzer/telemetry/logging"
)

func Handle(handler func(r *http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())

		response := handler(r)
		if response.Error != nil {
			err := response.Error
			logging.FromContext(r.Context()).Error("aplication err", zap.Error(err))
		}

		if h := response.Header(); h != nil {
			copyHeaders(w, h)
		}

		err := sendJSON(w, response.Body, response.Status)
		if err != nil {
			logging.FromContext(r.Context()).Error("on send json", zap.Error(err))
			span.RecordError(err)
		}
	}
}

func sendJSON(w http.ResponseWriter, payload any, statusCode int) error {
	if payload == nil {
		w.WriteHeader(statusCode)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(payload)
}

func copyHeaders(w http.ResponseWriter, h http.Header) {
	wh := w.Header()
	for header, values := range h {
		for _, value := range values {
			wh.Add(header, value)
		}
	}
}

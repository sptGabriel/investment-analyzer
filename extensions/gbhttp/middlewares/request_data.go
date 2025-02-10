package middlewares

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sptGabriel/investment-analyzer/extensions/utils"
)

func RequestDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, port := readUserIP(r)

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := utils.WithRequestData(r.Context(), utils.RequestData{
			IP:        ip,
			Port:      port,
			RequestID: requestID,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func readUserIP(r *http.Request) (string, string) {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}

	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	ip := "unknown"
	port := "unknown"
	if parts := strings.Split(IPAddress, ":"); len(parts) == 2 {
		ip = parts[0]
		port = parts[1]
	}
	return ip, port
}

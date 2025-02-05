package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
)

// Logger middleware for access logs
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Gathers metrics from the upstream handlers
			metrics := httpsnoop.CaptureMetrics(h, w, r)

			// Prints log and metrics
			logger.Info(
				strings.Join([]string{
					"handled request",
					slog.String("method", r.Method).Value.String(),
					slog.String("uri", r.URL.RequestURI()).Value.String(),
					slog.String("user_agent", r.Header.Get("User-Agent")).Value.String(),
					slog.String("ip", r.RemoteAddr).Value.String(),
					slog.Int("code", metrics.Code).Value.String(),
					slog.Int64("bytes", metrics.Written).Value.String(),
					slog.Duration("request_time", metrics.Duration).Value.String(),
				}, ""),
			)
		})
	}
}

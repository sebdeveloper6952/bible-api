package api

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder wraps ResponseWriter to capture the written status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(rec, r)

			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Int("status", rec.status),
				slog.Duration("latency", time.Since(start)),
			}

			switch {
			case rec.status >= 500:
				logger.Error("request", attrs...)
			case rec.status >= 400:
				logger.Warn("request", attrs...)
			default:
				logger.Info("request", attrs...)
			}
		})
	}
}

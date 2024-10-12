package middleware

import (
	"github.com/ankodd/url-shortener/internal/metrics"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func Logging(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(
				"Incoming request",
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
				slog.String("user-agent", r.UserAgent()),
				slog.String("remote-addr", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)
		})
	}
}

func Metrics(m *metrics.Metrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.IncRequest()
			next.ServeHTTP(w, r)
		})
	}
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

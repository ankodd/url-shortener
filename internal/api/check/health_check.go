package check

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HealthCheck handler for health check service
func HealthCheck(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.WithGroup("check.HealthCheck")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})

		log.Info("Health checked", slog.Bool("ok", true))
	}
}

package logger

import (
	"github.com/golang-cz/devslog"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

// Setup configured logger by env
//
// Uses github.com/golang-cz/devslog. Only on local env. For pretty logging
//
// Returning *slog.Logger
func Setup(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		opts := &devslog.Options{
			HandlerOptions:    &slog.HandlerOptions{Level: slog.LevelDebug},
			MaxSlicePrintSize: 10,
			SortKeys:          false,
			NewLineAfterLog:   true,
			StringerFormatter: true,
		}
		log = slog.New(devslog.NewHandler(os.Stdout, opts))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}

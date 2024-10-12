package main

import (
	"github.com/ankodd/url-shortener/internal/api/check"
	"github.com/ankodd/url-shortener/internal/api/url/redirect"
	"github.com/ankodd/url-shortener/internal/api/url/save"
	"github.com/ankodd/url-shortener/internal/config"
	"github.com/ankodd/url-shortener/internal/metrics"
	"github.com/ankodd/url-shortener/internal/middleware"
	"github.com/ankodd/url-shortener/internal/storage/postgres"
	"github.com/ankodd/url-shortener/pkg/logger"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Setup logger
	log := logger.Setup(cfg.Env)
	log.Debug("Logger configured")

	// Setup storage
	storage, err := postgres.New(&cfg.PostgreSQL)
	if err != nil {
		log.Error("Failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer storage.Close()
	log.Debug("Storage configured")

	// Setup metrics
	metric := metrics.NewMetrics()
	// initial router
	r := mux.NewRouter()
	r.Use(middleware.Logging(log), middleware.ContentTypeJSON, middleware.Metrics(metric))

	// Initial endpoints
	r.Handle("/save", save.Save(storage, log)).Methods(http.MethodPost)
	r.Handle("/{alias}", redirect.Redirect(storage, log)).Methods(http.MethodGet)
	r.Handle("/health-check", check.HealthCheck(log)).Methods(http.MethodGet)

	// Initial Server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      r,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
	log.Debug("Server configured")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := metrics.MostStartMetrics(cfg.MetricsAddr)
		if err != nil {
			log.Error("Failed start metrics server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	log.Debug("Metrics server listening", slog.String("address", cfg.MetricsAddr))

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = srv.ListenAndServe()
		if err != nil {
			log.Error("Failed start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	log.Debug("Server listening", slog.String("address", cfg.HTTPServer.Addr))
	wg.Wait()
}

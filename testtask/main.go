package main

import (
	"EMtask/testtask/adapters/api"
	"EMtask/testtask/adapters/db"
	"EMtask/testtask/adapters/rest"
	"EMtask/testtask/config"
	"EMtask/testtask/core"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "EMtask/testtask/adapters/rest/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Effective Mobile test task
// @description API для управления данными людей с обогащением информации через внешние API

// @host localhost:8080
// @BasePath /
func main() {
	// init config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// init logger
	log := mustMakeLogger(cfg.LogLevel)
	log.Info("starting server")

	// init api
	apiClient := api.NewClient(cfg, log)

	// init storage
	storage, err := db.New(log, cfg.DBAddress)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}

	// run migrations
	if err := storage.Migrate(); err != nil {
		log.Error("failed to migrate db", "error", err)
		os.Exit(1)
	}

	service := core.NewService(storage, apiClient, log)
	// init handlers
	handlers := rest.NewHandlers(service, log)

	mux := http.NewServeMux()

	// CRUD
	mux.HandleFunc("POST /people", handlers.CreatePersonHandler)
	mux.HandleFunc("GET /people", handlers.GetPeopleHandler)
	mux.HandleFunc("GET /person", handlers.GetPersonHandler)
	mux.HandleFunc("PUT /person", handlers.UpdatePersonHandler)
	mux.HandleFunc("DELETE /person", handlers.DeletePersonHandler)

	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	server := http.Server{
		Addr:        cfg.HTTPConfig.Address,
		Handler:     mux,
		ReadTimeout: cfg.HTTPConfig.Timeout,
	}

	log.Info("starting server", "address", cfg.HTTPConfig.Address)
	if err := server.ListenAndServe(); err != nil {
		log.Error("server error", "error", err)
	}

}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

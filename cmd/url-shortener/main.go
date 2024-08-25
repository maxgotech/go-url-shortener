package main

import (
	"log"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	}

	// load config
	cfg := config.MustLoad()

	// create logger
	log := setupLogger((cfg.Env))

	log.Info(("starting url-shortener"), slog.String("env", cfg.Env))

	log.Debug("debug messages enabled")

	// create storage
	storage, err := sqlite.NewStorage(cfg.StoragePath, log)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// create router
	router := chi.NewRouter()

	// id for each req
	router.Use(middleware.RequestID)
	// ip of user for req
	router.Use(middleware.RealIP)

	_ = router
	// TODO: run server:
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

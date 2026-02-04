package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/resume-builder/backend/internal/api"
	"github.com/resume-builder/backend/internal/config"
	"github.com/resume-builder/backend/internal/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level based on environment
	if cfg.IsProduction() {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().
		Str("environment", cfg.Environment).
		Str("port", cfg.Port).
		Msg("Starting Resume Builder API")

	// Connect to database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Setup router
	router := api.NewRouter(cfg, db)
	app := router.Setup()

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info().Msg("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Server shutdown error")
		}
	}()

	// Start server
	addr := ":" + cfg.Port
	log.Info().Str("addr", addr).Msg("Server listening")
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

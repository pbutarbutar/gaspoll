package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gaspoll/internal/config"
	"gaspoll/internal/db"
	"gaspoll/internal/observability"
	"gaspoll/internal/server"
)

func main() {
	setupLogger()

	cfg := config.Load()

	shutdownTelemetry := observability.Setup(cfg)
	defer shutdownTelemetry(context.Background())

	database := db.MustConnect(cfg)
	defer database.Close()

	srv, err := server.New(cfg, database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init server")
	}

	go func() {
		addr := ":" + cfg.HTTPPort
		log.Info().Str("addr", addr).Msg("server listening")
		if err := srv.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
	log.Info().Msg("server stopped")
}

func setupLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

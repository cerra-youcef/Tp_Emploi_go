package main

import (
	"alerter/internal/alerter"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ConfigURL    string
	TimetableURL string
}

func loadConfig() (Config, error) {
	var cfg Config
	var ok bool

	if cfg.ConfigURL, ok = os.LookupEnv("CONFIG_URL"); !ok {
		return cfg, errors.New("CONFIG_URL not set in .env file")
	}

	if cfg.TimetableURL, ok = os.LookupEnv("TIMETABLE_URL"); !ok {
		return cfg, errors.New("TIMETABLE_URL not set in .env file")
	}

	return cfg, nil
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		slog.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	// Load configuration
	_, err := loadConfig()
	if err != nil {
		slog.Error("Configuration error", "error", err)
		os.Exit(1)
	}

	// Initialize NATS
	if err := alerter.InitNATS(); err != nil {
		slog.Error("Failed to initialize NATS", "error", err)
		os.Exit(1)
	}

	// Start message consumer
	go func() {
		slog.Info("Starting message consumer")
		if err := alerter.ConsumeMessages(ctx); err != nil {
			slog.Error("Message consumer failed", "error", err)
			cancel() // Trigger shutdown if consumer fails
		}
	}()

	// Wait for shutdown
	<-ctx.Done()

	// Allow a brief period for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	slog.Info("Shutting down...")

	<-shutdownCtx.Done()
	slog.Info("Shutdown complete")
}

package main

import (
	"alerter/internal/alerter"
	"alerter/internal/mailer"
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
	mailToken    string
	apiURL       string
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

	if cfg.mailToken, ok = os.LookupEnv("MAIL_TOKEN"); !ok {
		return cfg, errors.New("MAIL_TOKEN not set in .env file")
	}

	if cfg.apiURL, ok = os.LookupEnv("API_URL"); !ok {
		return cfg, errors.New("MAIL_TOKEN not set in .env file")
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
	cfg, err := loadConfig()
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

	// Exemple d'utilisation
	to := "yanis.beldjilali@etu.uca.fr"
	data := mailer.TemplateData{
		EventName:   "Réunion de projet",
		NewDate:     "15 mars 2024",
		NewLocation: "Salle B203",
	}

	// Générer le contenu de l'email
	subject, content, err := mailer.GetEmailContent("templates/mail.html", data)
	if err != nil {
		slog.Error("error while generation email template" + err.Error())
		return
	}

	err = mailer.SendEmail(to, subject, content, cfg.mailToken, cfg.apiURL)
	if err != nil {
		slog.Error("MAIL ERROR", "error", err)
	} else {
		slog.Error("NO ERROR", "succes", err)
	}

	// Wait for shutdown
	<-ctx.Done()

	// Allow a brief period for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	slog.Info("Shutting down...")

	<-shutdownCtx.Done()
	slog.Info("Shutdown complete")
}

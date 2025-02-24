package main

import (
	"log"
	"scheduler/internal/nats"
	"scheduler/internal/scheduler"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Scheduler Service Started")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configURL := os.Getenv("CONFIG_URL")
	if configURL == "" {
		log.Fatal("CONFIG_URL not set in .env file")
	}

	ucaURL := os.Getenv("UCA_URL")
	if ucaURL == "" {
		log.Fatal("UCA_URL not set in .env file")
	}

	// Initialize NATS connection
	nats.InitNATS()

	// Start scheduler
	scheduler.StartScheduler(configURL, ucaURL)
}

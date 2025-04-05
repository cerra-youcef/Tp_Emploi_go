package main

import (
	"alerter/internal/alerter"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configURL := os.Getenv("CONFIG_URL")
	if configURL == "" {
		log.Fatal("CONFIG_URL not set in .env file")
	}

	ttURL := os.Getenv("TIMETABLE_URL")
	if ttURL == "" {
		log.Fatal("TIMETABLE_URL not set in .env file")
	}

	//Start NATS Consumer in a Goroutine with context
	go func() {
		log.Println("Starting NATS Consumer...")
		js, nc, err := alerter.ConnectToNATS()
		if err != nil {
			log.Printf("Error connecting to NATS: %v", err)
			return //to keep the api runnig
		}
		defer nc.Close()

		consumer, err := alerter.AlertConsumer(js)
		if err != nil {
			log.Printf("Error creating NATS Consumer: %v", err)
			return //to keep the api runnig
		}

		err = alerter.Consume(*consumer) // Pass the context with DB
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
			return //to keep the api runnig
		}
	}()
}

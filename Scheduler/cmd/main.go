package main

import (
	"log"
	"scheduler/internal/nats"
	"scheduler/internal/scheduler"
)

func main() {
	log.Println("Scheduler Service Started")

	// Initialize NATS connection
	nats.InitNATS()

	// Start scheduler
	scheduler.StartScheduler()
}

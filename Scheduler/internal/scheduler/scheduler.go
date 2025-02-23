package scheduler

import (
	"context"
	"log"
	"time"

	"scheduler/internal/apiClients"
	"scheduler/internal/nats"

	"github.com/zhashkevych/scheduler"
)

func FetchAndPublishEvents(ctx context.Context) {
	configURL := "http://localhost:8080"

	// Fetch resources
	resources, err := apiClients.FetchResourcesFromConfig(configURL)
	if err != nil {
		log.Println("Error fetching resources:", err)
		return
	}

	// Extract resource IDs
	var resourceIDs []int
	for _, r := range resources {
		resourceIDs = append(resourceIDs, r.UcaId)
	}

	// Fetch events for these resources
	events, err := apiClients.FetchEventsFromUCA(resourceIDs)
	if err != nil {
		log.Println("Error fetching events:", err)
		return
	}

	// Publish events to NATS
	for _, event := range events {
		err := nats.PublishEvent("EVENTS.create", event)
		if err != nil {
			log.Println("Error publishing event:", err)
		}
	}

	log.Println("Successfully fetched & published events")
}

// Start scheduler
func StartScheduler() {
	ctx := context.Background()
	sc := scheduler.NewScheduler()
	sc.Add(ctx, FetchAndPublishEvents, time.Second*5) // Runs every 1 mins

	// Keep program alive
	quit := make(chan struct{})
	<-quit
	sc.Stop()
}

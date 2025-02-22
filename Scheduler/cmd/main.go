package main

import (
	"encoding/json"
	"log"
	"fmt"
	"scheduler/internal/apiClients"
)

func main() {
	configURL := "http://localhost:8080" 
	
	/*
	resourceIDs := []int{13295, 13345} // Example resource IDs
	events, err := apiClients.FetchEventsFromUCA(resourceIDs)
	if err != nil {
		fmt.Println("Error fetching events:", err)
		return
	}*/


	timetables, err := apiClients.FetchResourcesFromConfig(configURL)
	if err != nil {
		log.Fatalf("Error fetching timetables: %v", err)
	}

	// Convert to JSON and print
	jsonData, _ := json.MarshalIndent(timetables, "", "  ")
	fmt.Println(string(jsonData))
}

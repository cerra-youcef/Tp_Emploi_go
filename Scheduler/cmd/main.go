package main

import (
	"encoding/json"
	"log"
	"scheduler/internal/edt"
)

func main() {
	// Retrieve data from EDT
	events, err := edt.FetchEvents("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=13295,13345&projectId=2&calType=ical&nbWeeks=8&displayConfigId=128")
	if err != nil {
		log.Fatalf("Failed to fetch events: %v", err)
	}

	// here you need to do a call request tp an  endpoint
	//post(events)

	// Parse to JSON and display
	jsonData, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal events to JSO bN: %v", err)
	}

	log.Println(string(jsonData))
}

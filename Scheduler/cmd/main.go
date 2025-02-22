package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"fmt"
	"io"
	"net/http"
	"strings"
	"scheduler/internal/models"
)

// FetchEventsFromUCA fetches events based on provided resource IDs.
func FetchEventsFromUCA(resourceIDs []int) ([]models.Event, error) {
	if len(resourceIDs) == 0 {
		return nil, fmt.Errorf("no resource IDs provided")
	}

	// Convert resourceIDs to a comma-separated string
	resourceStr := strings.Trim(strings.Replace(fmt.Sprint(resourceIDs), " ", ",", -1), "[]")

	// Construct the URL dynamically
	url := fmt.Sprintf("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%s&projectId=2&calType=ical&nbWeeks=8&displayConfigId=128", resourceStr)

	// Retrieve data from EDT
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	// Read all data from response
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Create a line-reader from data
	scanner := bufio.NewScanner(bytes.NewReader(rawData))

	// Store parsed events
	var events []models.Event
	currentEvent := models.Event{}

	currentKey := ""
	currentValue := ""
	inEvent := false

	// Parse each line
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore non-event lines
		if !inEvent && line != "BEGIN:VEVENT" {
			continue
		}

		// Start a new event
		if line == "BEGIN:VEVENT" {
			inEvent = true
			currentEvent = models.Event{}
			continue
		}

		// End event and store
		if line == "END:VEVENT" {
			inEvent = false
			events = append(events, currentEvent)
			continue
		}

		// Handle multi-line values (continuation lines start with a space)
		if strings.HasPrefix(line, " ") {
			currentValue += strings.TrimSpace(line)
			currentEvent = updateEventField(currentEvent, currentKey, currentValue)
			continue
		}

		// Split key-value pair
		splitted := strings.SplitN(line, ":", 2)
		if len(splitted) < 2 {
			continue
		}
		currentKey = splitted[0]
		currentValue = splitted[1]

		// Store field
		currentEvent = updateEventField(currentEvent, currentKey, currentValue)
	}

	return events, nil
}

func updateEventField(event models.Event, key, value string) models.Event {
	switch key {
	case "UID":
		event.UID = value
	case "SUMMARY":
		event.Name = value
	case "DESCRIPTION":
		event.Description = value
	case "LOCATION":
		event.Location = value
	case "DTSTART":
		event.Start = value
	case "DTEND":
		event.End = value
	case "CREATED":
		event.CreatedAt = value
	case "LAST-MODIFIED":
		event.UpdatedAt = value
	case "DTSTAMP":
		event.DTStamp = value
	}
	return event
}

func fetchResourcesFromConfig(configURL string) ([]models.Resource, error) {
	resp, err := http.Get(fmt.Sprintf("%s/resources", configURL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %w", err)
	}
	defer resp.Body.Close()

	// Decode JSON response into a generic slice of maps
	var rawResources []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResources); err != nil {
		return nil, fmt.Errorf("failed to decode resources response: %w", err)
	}

	// Convert the raw data into models.Resource
	var timetables []models.Resource
	for _, raw := range rawResources {
		idStr, ok := raw["id"].(string)
		if !ok {
			log.Println("Invalid or missing 'id' field")
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Failed to parse UUID for resource ID %s: %v", idStr, err)
			continue
		}

		ucaId, ok := raw["uca_id"].(float64) // JSON numbers are float64 by default
		if !ok {
			log.Println("Invalid or missing 'uca_id' field")
			continue
		}

		name, ok := raw["name"].(string)
		if !ok {
			log.Println("Invalid or missing 'name' field")
			continue
		}

		timetables = append(timetables, models.Resource{
			ID:    id,
			UcaId: int(ucaId), // Convert float64 to int
			Name:  name,
		})
	}

	return timetables, nil
}


func main() {
	/*
	resourceIDs := []int{13295, 13345} // Example resource IDs
	events, err := FetchEventsFromUCA(resourceIDs)
	if err != nil {
		fmt.Println("Error fetching events:", err)
		return
	}*/

	configURL := "http://localhost:8080" 

	timetables, err := fetchResourcesFromConfig(configURL)
	if err != nil {
		log.Fatalf("Error fetching timetables: %v", err)
	}

	// Convert to JSON and print
	jsonData, _ := json.MarshalIndent(timetables, "", "  ")
	fmt.Println(string(jsonData))
}

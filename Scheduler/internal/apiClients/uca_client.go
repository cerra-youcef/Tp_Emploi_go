package apiClients

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"scheduler/internal/models"
)

// FetchEventsFromUCA iterates over resource IDs and fetches events separately for each.
func FetchEventsFromUCA(resourceIDs []int) ([]models.Event, error) {
	if len(resourceIDs) == 0 {
		return nil, fmt.Errorf("no resource IDs provided")
	}

	var allEvents []models.Event

	for _, resourceID := range resourceIDs {
		events, err := fetchEventsForSingleResource(resourceID)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for resource %d: %v", resourceID, err)
		}

		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}

// fetchEventsForSingleResource fetches events for a single resource ID
func fetchEventsForSingleResource(resourceID int) ([]models.Event, error) {
	url := buildUCAURL(resourceID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return parseEvents(rawData, resourceID), nil
}

// buildUCAURL constructs the UCA API URL for a given resource ID
func buildUCAURL(resourceID int) string {
	return fmt.Sprintf(
		"https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%d&projectId=2&calType=ical&nbWeeks=8&displayConfigId=128",
		resourceID,
	)
}

// parseEvents parses the raw iCalendar data and assigns the correct resource ID
func parseEvents(rawData []byte, resourceID int) []models.Event {
	scanner := bufio.NewScanner(bytes.NewReader(rawData))

	var events []models.Event
	currentEvent := models.Event{}
	currentKey := ""
	currentValue := ""
	inEvent := false

	for scanner.Scan() {
		line := scanner.Text()

		if !inEvent && line != "BEGIN:VEVENT" {
			continue
		}

		if line == "BEGIN:VEVENT" {
			inEvent = true
			currentEvent = models.Event{
				Resource: resourceID,
			}
			continue
		}

		if line == "END:VEVENT" {
			inEvent = false
			events = append(events, currentEvent)
			continue
		}

		if strings.HasPrefix(line, " ") {
			currentValue += strings.TrimSpace(line)
			currentEvent = updateEventField(currentEvent, currentKey, currentValue)
			continue
		}

		splitted := strings.SplitN(line, ":", 2)
		if len(splitted) < 2 {
			continue
		}
		currentKey = splitted[0]
		currentValue = splitted[1]

		currentEvent = updateEventField(currentEvent, currentKey, currentValue)
	}

	return events
}

// updateEventField maps iCalendar fields to the Event struct
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

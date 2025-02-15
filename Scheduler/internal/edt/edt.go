package edt

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid" // <-- Import du package uuid
	"scheduler/internal/models"
)

func FetchEvents(url string) ([]models.Event, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(rawData))
	var events []models.Event
	var currentEvent models.Event
	inEvent := false

	for scanner.Scan() {
		line := scanner.Text()

		if !inEvent && line != "BEGIN:VEVENT" {
			continue
		}

		if line == "BEGIN:VEVENT" {
			inEvent = true
			currentEvent = models.Event{}
			continue
		}

		if line == "END:VEVENT" {
			// Generate a new UUID for the event
			id := uuid.New()
			currentEvent.ID = id

			// Append the event to the list
			events = append(events, currentEvent)
			inEvent = false
			continue
		}

		if strings.HasPrefix(line, " ") {
			// Handle multi-line data (if needed)
			continue
		}

		splitted := strings.SplitN(line, ":", 2)
		if len(splitted) != 2 {
			continue
		}

		key := splitted[0]
		value := splitted[1]

		switch key {
		case "UID":
			currentEvent.UID = value
		case "DESCRIPTION":
			currentEvent.Description = value
		case "SUMMARY":
			currentEvent.Name = value
		case "DTSTART":
			startTime, err := time.Parse("20060102T150405Z", value)
			if err == nil {
				currentEvent.Start = startTime
			}
		case "DTEND":
			currentEvent.End = value
		case "LOCATION":
			currentEvent.Location = value
		case "LAST-MODIFIED":
			lastUpdate, err := time.Parse("20060102T150405Z", value)
			if err == nil {
				currentEvent.UpdatedAt = lastUpdate
			}
		case "RESOURCES":
			// Handle resource IDs (assuming they are comma-separated UUIDs)
			resourceIDs := strings.Split(value, ",")
			for _, idStr := range resourceIDs {
				id, err := uuid.Parse(idStr)
				if err == nil {
					currentEvent.RessourceIDs = append(currentEvent.RessourceIDs, id)
				}
			}
		}
	}

	return events, nil
}
package edt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"scheduler/internal/models"
)

// FetchEventsFromUCA récupère les événements pour un emploi du temps spécifique.
func FetchEventsFromUCA(resourceID int) ([]models.Event, error) {
	// Construire l'URL pour l'API UCA
	url := fmt.Sprintf("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%d&projectId=2&calType=ical&nbWeeks=8&displayConfigId=128", resourceID)

	// Effectuer la requête HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events from UCA: %w", err)
	}
	defer resp.Body.Close()

	// Lire le contenu du fichier .ics
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Scanner le fichier ligne par ligne
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
			// Générer un UUID unique pour l'événement
			id := uuid.New()
			currentEvent.ID = id
			events = append(events, currentEvent)
			inEvent = false
			continue
		}

		// Ignorer les lignes multi-lignes
		if strings.HasPrefix(line, " ") {
			continue
		}

		// Analyser les champs clé-valeur
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
		}
	}

	return events, nil
}
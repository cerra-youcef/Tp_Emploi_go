package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"scheduler/internal/models"
	"scheduler/internal/edt"
)

// Structure temporaire pour déserialiser la réponse JSON de l'API Config
type ResourceResponse struct {
	ID    string `json:"id"`    // ID en format chaîne (converti en UUID plus tard)
	UcaId int    `json:"uca_id"`
	Name  string `json:"name"`
}

// Fonction pour récupérer les emplois du temps depuis l'API "Config"
func fetchTimetablesFromConfig(configURL string) ([]models.Ressource, error) {
	resp, err := http.Get(fmt.Sprintf("%s/resources", configURL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %w", err)
	}
	defer resp.Body.Close()

	var resources []ResourceResponse
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("failed to decode resources response: %w", err)
	}

	// Convertir les données en modèle interne
	var timetables []models.Ressource
	for _, r := range resources {
		id, err := uuid.Parse(r.ID)
		if err != nil {
			log.Printf("Failed to parse UUID for resource ID %s: %v", r.ID, err)
			continue
		}
		timetables = append(timetables, models.Ressource{
			ID:    id,
			UcaId: r.UcaId,
			Name:  r.Name,
		})
	}

	return timetables, nil
}

// Génère une représentation iCalendar d'un événement
func generateICalEvent(event models.Event) string {
	return fmt.Sprintf(
		"BEGIN:VEVENT\n"+
			"DTSTAMP:%s\n"+
			"DTSTART:%s\n"+
			"DTEND:%s\n"+
			"SUMMARY:%s\n"+
			"LOCATION:%s\n"+
			"DESCRIPTION:%s\n"+
			"UID:%s\n"+
			"LAST-MODIFIED:%s\n"+
			"END:VEVENT\n",
		event.Start.Format("20060102T150405Z"),
		event.Start.Format("20060102T150405Z"),
		event.End,
		event.Name,
		event.Location,
		strings.ReplaceAll(event.Description, "\n", "\\n"),
		event.UID,
		event.UpdatedAt.Format("20060102T150405Z"),
	)
}

// Génère une représentation complète iCalendar pour une liste d'événements
func generateICal(events []models.Event) string {
	var ical strings.Builder
	ical.WriteString("BEGIN:VCALENDAR\n")
	ical.WriteString("METHOD:PUBLISH\n")
	ical.WriteString("PRODID:-//ADE/version 6.0\n")
	ical.WriteString("VERSION:2.0\n")
	ical.WriteString("CALSCALE:GREGORIAN\n")

	for _, event := range events {
		ical.WriteString(generateICalEvent(event))
	}

	ical.WriteString("END:VCALENDAR\n")
	return ical.String()
}

// Simule l'envoi des événements vers une MQ
func sendEventsToMQ(events []models.Event) {
	for _, event := range events {
		fmt.Printf("Sending event to MQ: %s (%s)\n", event.Name, event.ID)
	}
}

func main() {
	configURL := "http://localhost:8080" // URL de l'API Config

	// Étape 1 : Récupérer les emplois du temps configurés depuis l'API "Config"
	timetables, err := fetchTimetablesFromConfig(configURL)
	if err != nil {
		log.Fatalf("Error fetching timetables: %v", err)
	}

	// Afficher les emplois du temps récupérés
	fmt.Println("Fetched Timetables:")
	for _, t := range timetables {
		fmt.Printf("ID: %s, UCA ID: %d, Name: %s\n", t.ID, t.UcaId, t.Name)
	}

	// Étape 2 : Pour chaque emploi du temps, récupérer les événements depuis l'API UCA
	for _, timetable := range timetables {
		fmt.Printf("\nFetching events for timetable: %s (UCA ID: %d)\n", timetable.Name, timetable.UcaId)

		// Récupérer les événements depuis l'API UCA
		events, err := edt.FetchEventsFromUCA(timetable.UcaId)
		if err != nil {
			log.Printf("Error fetching events for timetable %s: %v", timetable.Name, err)
			continue
		}

		// Afficher les événements récupérés
		fmt.Printf("Fetched %d events for timetable %s:\n", len(events), timetable.Name)
		for _, event := range events {
			fmt.Printf("- Event: %s, Start: %s, End: %s, Location: %s\n",
				event.Name, event.Start.Format(time.RFC3339), event.End, event.Location)
		}

		// Générer et afficher la représentation iCalendar
		icalContent := generateICal(events)
		fmt.Printf("\nGenerated iCalendar for timetable %s:\n%s\n", timetable.Name, icalContent)

		// Envoyer les événements vers la MQ
		sendEventsToMQ(events)
	}

	fmt.Println("\nScheduler execution completed.")
}
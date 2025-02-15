package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"scheduler/internal/models"
)

// Structure temporaire pour déserialiser la réponse JSON
type ResourceResponse struct {
	ID    uuid.UUID `json:"id"`
	UcaId int       `json:"uca_id"`
	Name  string    `json:"name"`
}

// Fonction pour récupérer toutes les ressources depuis l'API "Config"
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
		timetables = append(timetables, models.Ressource{
			ID:    r.ID,
			UcaId: r.UcaId,
			Name:  r.Name,
		})
	}

	return timetables, nil
}

func main() {
	configURL := "http://localhost:8080" // URL de l'API Config

	// Récupérer les emplois du temps configurés
	timetables, err := fetchTimetablesFromConfig(configURL)
	if err != nil {
		log.Fatalf("Error fetching timetables: %v", err)
	}

	// Afficher les emplois du temps récupérés
	fmt.Println("Fetched Timetables:")
	for _, t := range timetables {
		fmt.Printf("ID: %s, UCA ID: %d, Name: %s\n", t.ID, t.UcaId, t.Name)
	}
}
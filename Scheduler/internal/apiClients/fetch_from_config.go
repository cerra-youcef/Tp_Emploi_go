package apiClients

import (
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"fmt"
	"net/http"
	"scheduler/internal/models"
)

func FetchResourcesFromConfig(configURL string) ([]models.Resource, error) {
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

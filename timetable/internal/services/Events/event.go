package Events

import (
	"github.com/google/uuid"
	"timetable/internal/models"
	"timetable/internal/repositories/Events"
)

// GetEventsByResourceID récupère les événements associés à une ressource spécifique.
func GetEventsByResourceID(resourceID string) ([]models.Event, error) {
	return Events.GetEventsByResourceID(resourceID)
}

// GetAllEvents récupère tous les événements.
func GetAllEvents() ([]models.Event, error) {
	return Events.GetAllEvents()
}

// GetEventByID récupère un événement par son ID.
func GetEventByID(id uuid.UUID) (*models.Event, error) {
	return Events.GetEventByID(id)
}

// CreateEvent crée un nouvel événement.
func CreateEvent(event *models.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	return Events.CreateEvent(event)
}
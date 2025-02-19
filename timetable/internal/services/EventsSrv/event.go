package EventsSrv

import (
"github.com/google/uuid"
"timetable/internal/models"
"timetable/internal/repositories/Eventsrep"
"net/http" // Ajoutez cet import si nécessaire
"database/sql"


)

// GetEventsByResourceID récupère les événements associés à une ressource spécifique.
func GetEventsByResourceID(db *sql.DB, resourceID string) ([]*models.Event, error) {
	return Eventsrep.GetEventsByResourceID(db, resourceID)
}

// GetAllEvents récupère tous les événements.
func GetAllEvents() ([]models.Event, error) {
	return Eventsrep.GetAllEvents()
}

// GetEventByID récupère un événement par son ID.
func GetEventByID(id uuid.UUID) (*models.Event, error) {
	return Eventsrep.GetEventByID(id)
}

// CreateEvent crée un nouvel événement.
func CreateEvent(event *models.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	return Eventsrep.CreateEvent(event)
}

// UpdateEvent met à jour un événement existant.
func UpdateEvent(id uuid.UUID, event *models.Event) error {
	if event.ID != id {
		return &models.CustomError{
			Message: "ID mismatch between URL and request body",
			Code:    http.StatusBadRequest,
		}
	}
	return Eventsrep.UpdateEvent(id, event)
}

// DeleteEvent supprime un événement existant.
func DeleteEvent(id uuid.UUID) error {
	return Eventsrep.DeleteEvent(id)
}

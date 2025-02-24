package Events

import (
	"github.com/google/uuid"
	"timetable/internal/models"
	"timetable/internal/repositories/Events"
	"timetable/internal/helpers"
	"timetable/internal/nats/publisher"
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

func GetEventByUID(id string) (*models.Event, error) {
	return Events.GetEventByUID(id)
}

func UpdateEvent(event *models.Event) error {
	return Events.UpdateEvent(event)
}

func DeleteEvent(id uuid.UUID) error {
	return Events.DeleteEvent(id)
}


func CreateAndNotifyEvent(event models.Event) error {
	if err := Events.CreateEvent(&event); err != nil {
		return err
	}
	return natsPublisher.PublishAlert("create", helpers.CreateAlert("event.created", event))
}

func UpdateAndNotifyEvent(oldEvent, newEvent models.Event) error {
	changes := helpers.DetectChanges(oldEvent, newEvent)
	if len(changes) == 0 {
		return nil
	}
	if err := Events.UpdateEvent(&newEvent); err != nil {
		return err
	}
	return natsPublisher.PublishAlert("updated", helpers.CreateAlert("event.updated", newEvent, changes))
}
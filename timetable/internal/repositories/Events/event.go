package Events

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"timetable/internal/helpers"
	"timetable/internal/models"
)

// GetAllEvents récupère tous les événements.
func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, resource_ids, uid, name, start
        FROM events
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		var resourceIdsJSON string
		err = rows.Scan(&event.ID, &resourceIdsJSON, &event.UID, &event.Name, &event.Start)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(resourceIdsJSON), &event.ResourceIDs)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// GetEventByID récupère un événement par son ID.
func GetEventByID(id uuid.UUID) (*models.Event, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, resource_ids, uid, name, start
        FROM events
        WHERE id = ?
    `

	var event models.Event
	var resourceIdsJSON string
	err = db.QueryRow(query, id).Scan(&event.ID, &resourceIdsJSON, &event.UID, &event.Name, &event.Start)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retourne nil si l'événement n'existe pas.
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(resourceIdsJSON), &event.ResourceIDs)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// CreateEvent crée un nouvel événement.
func CreateEvent(event *models.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	resourceIdsJSON, err := json.Marshal(event.ResourceIDs)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO events (id, resource_ids, uid, name, start)
        VALUES (?, ?, ?, ?, ?)
    `

	result, err := db.Exec(query, event.ID, string(resourceIdsJSON), event.UID, event.Name, event.Start)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

func GetEventsByResourceID(db *sql.DB, resourceID string) ([]*models.Event, error) {
	rows, err := db.Query("SELECT id, resource_ids, uid, name, start, created_at, updated_at FROM events WHERE resource_ids LIKE ?", "%"+resourceID+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.ResourceIDs, &event.UID, &event.Name, &event.Start, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

// UpdateEvent met à jour un événement existant.
func UpdateEvent(id uuid.UUID, event *models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	resourceIdsJSON, err := json.Marshal(event.ResourceIDs)
	if err != nil {
		return err
	}

	query := `
        UPDATE events
        SET resource_ids = ?, uid = ?, name = ?, start = ?
        WHERE id = ?
    `

	result, err := db.Exec(query, string(resourceIdsJSON), event.UID, event.Name, event.Start, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// DeleteEvent supprime un événement existant.
func DeleteEvent(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        DELETE FROM events
        WHERE id = ?
    `

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}
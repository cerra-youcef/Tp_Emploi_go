package Alertsrep

import (
	"database/sql"
	"github.com/google/uuid"
	"cerra/tp_go/internal/helpers"
	"cerra/tp_go/internal/models"
)

// GetAllAlerts récupère toutes les alertes.
func GetAllAlerts() ([]models.Alert, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, email, is_all, ressource_id
        FROM alerts
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var alert models.Alert
		err = rows.Scan(&alert.ID, &alert.Email, &alert.All, &alert.RessourceID)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetAlertByID récupère une alerte par son ID.
func GetAlertByID(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, email, is_all, ressource_id
        FROM alerts
        WHERE id = ?
    `

	var alert models.Alert
	err = db.QueryRow(query, id).Scan(&alert.ID, &alert.Email, &alert.All, &alert.RessourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retourne nil si l'alerte n'existe pas.
		}
		return nil, err
	}

	return &alert, nil
}

// CreateAlert crée une nouvelle alerte.
func CreateAlert(alert *models.Alert) error {
	if alert.ID == uuid.Nil {
		alert.ID = uuid.New()
	}

	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        INSERT INTO alerts (id, email, is_all, ressource_id)
        VALUES (?, ?, ?, ?)
    `

	result, err := db.Exec(query, alert.ID, alert.Email, alert.All, alert.RessourceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// UpdateAlert met à jour une alerte existante.
func UpdateAlert(id uuid.UUID, alert *models.Alert) error {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        UPDATE alerts
        SET email = ?, is_all = ?, ressource_id = ?
        WHERE id = ?
    `

	result, err := db.Exec(query, alert.Email, alert.All, alert.RessourceID, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// DeleteAlert supprime une alerte existante.
func DeleteAlert(id uuid.UUID) error {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        DELETE FROM alerts
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
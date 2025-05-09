package Alerts

import (
	"config/internal/models"
	"config/internal/repositories/Alerts"
	"net/http" // Ajoutez cet import si nécessaire

	"github.com/google/uuid"
)

// GetAllAlerts récupère toutes les alertes.
func GetAllAlerts(ucaID string) ([]models.Alert, error) {
	return Alerts.GetAllAlerts(ucaID)
}

// GetAlertByID récupère une alerte par son ID.
func GetAlertByID(id uuid.UUID) (*models.Alert, error) {
	return Alerts.GetAlertByID(id)
}

// CreateAlert crée une nouvelle alerte.
func CreateAlert(alert *models.Alert) error {
	if alert.ID == uuid.Nil {
		alert.ID = uuid.New()
	}
	return Alerts.CreateAlert(alert)
}

// UpdateAlert met à jour une alerte existante.
func UpdateAlert(id uuid.UUID, alert *models.Alert) error {
	_, err := GetAlertByID(id)
	if err != nil {
		return &models.CustomError{
			Message: "Alert not found ",
			Code:    http.StatusNotFound,
		}
	}

	return Alerts.UpdateAlert(id, alert)
}

// DeleteAlert supprime une alerte existante.
func DeleteAlert(id uuid.UUID) error {
	return Alerts.DeleteAlert(id)
}

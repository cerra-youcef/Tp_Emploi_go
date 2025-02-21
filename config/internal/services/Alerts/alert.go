package Alerts

import (
	"github.com/google/uuid"
	"config/internal/models"
	"config/internal/repositories/Alerts"
	"net/http" // Ajoutez cet import si nécessaire

)

// GetAllAlerts récupère toutes les alertes.
func GetAllAlerts() ([]models.Alert, error) {
	return Alerts.GetAllAlerts()
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
	if alert.ID != id {
		return &models.CustomError{
			Message: "ID mismatch between URL and request body",
			Code:    http.StatusBadRequest,
		}
	}
	return Alerts.UpdateAlert(id, alert)
}

// DeleteAlert supprime une alerte existante.
func DeleteAlert(id uuid.UUID) error {
	return Alerts.DeleteAlert(id)
}
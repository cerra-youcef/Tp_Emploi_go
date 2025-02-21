package Alerts

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"config/internal/models"
	"config/internal/services/Alerts"
	"github.com/google/uuid"


	"net/http"
)

// UpdateAlertHandler updates an existing alert.
// @Summary Update an alert by its ID
// @Description This endpoint allows you to update an alert by providing the alert ID and the fields to be updated.
// @Tags Alerts
// @Accept json
// @Produce json
// @Param alertId path string true "Alert ID" example("123e4567-e89b-12d3-a456-426614174000")
// @Param alert body models.Alert true "Updated alert object"
// @Success 204 "Alert updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Alert not found"
// @Failure 500 {string} string "Internal server error"
// @Router /alerts/{alertId} [put]
func UpdateAlertHandler(w http.ResponseWriter, r *http.Request) {
	alertIdRaw := r.Context().Value("alertId") 
	if alertIdRaw == nil {
		logrus.Errorf("Alert ID not found in context")
		http.Error(w, "Missing alert in context", http.StatusBadRequest)
		return
	}

	alertId, ok := alertIdRaw.(uuid.UUID)
	if !ok {
		logrus.Errorf("Alert ID is not of type uuid.UUID")
		http.Error(w, "Invalid alert ID type", http.StatusBadRequest)
		return
	}

	var updatedAlert models.Alert
	err := json.NewDecoder(r.Body).Decode(&updatedAlert)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = Alerts.UpdateAlert(alertId, &updatedAlert)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error updating alert: %v", err)
		http.Error(w, "Failed to update alert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

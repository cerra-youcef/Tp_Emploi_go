package Alerts

import (
"encoding/json"
"github.com/sirupsen/logrus"
"config/internal/models"
"config/internal/services/Alerts"

"net/http"
)

// CreateAlertHandler creates a new alert.
// @Summary Create a new alert
// @Description This endpoint allows you to create a new alert by providing the necessary details.
// @Tags Alerts
// @Accept json
// @Produce json
// @Param alert body models.Alert true "Alert object"
// @Success 201 {object} models.Alert "The newly created alert"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /alerts [post]
func CreateAlertHandler(w http.ResponseWriter, r *http.Request) {
	var newAlert models.Alert
	err := json.NewDecoder(r.Body).Decode(&newAlert)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = Alerts.CreateAlert(&newAlert)
	if err != nil {
		logrus.Errorf("Error creating alert: %v", err)
		http.Error(w, "Failed to create alert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAlert)
}

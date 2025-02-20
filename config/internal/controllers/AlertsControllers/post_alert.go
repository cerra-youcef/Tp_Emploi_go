package AlertsControllers

import (
"encoding/json"
"github.com/sirupsen/logrus"
"config/internal/models"
"config/internal/services/AlertsSrv"

"net/http"
)

// CreateAlertHandler gère la création d'une nouvelle alerte.
func CreateAlertHandler(w http.ResponseWriter, r *http.Request) {
	var newAlert models.Alert
	err := json.NewDecoder(r.Body).Decode(&newAlert)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = AlertsSrv.CreateAlert(&newAlert)
	if err != nil {
		logrus.Errorf("Error creating alert: %v", err)
		http.Error(w, "Failed to create alert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAlert)
}

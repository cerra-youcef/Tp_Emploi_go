package RessourcesControllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"config/internal/models"
	"config/internal/services/RessourcesSrv"
	"net/http"
	"github.com/google/uuid"

)

// CreateResourceHandler crée une nouvelle ressource.
func CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	var newResource models.Ressource
	err := json.NewDecoder(r.Body).Decode(&newResource)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newResource.ID == uuid.Nil {
		newResource.ID = uuid.New()
	}

	err = RessourcesSrv.CreateResource(&newResource)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error creating resource: %v", err)
		http.Error(w, "Failed to create resource", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newResource)
}


package RessourcesControllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"cerra/tp_go/internal/models"
	"cerra/tp_go/internal/services/RessourcesSrv"
	"net/http"
	"github.com/google/uuid"

)

// UpdateResourceHandler met à jour une ressource existante.
func UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
	resourceId := r.Context().Value("resourceId").(uuid.UUID)

	var updatedResource models.Ressource
	err := json.NewDecoder(r.Body).Decode(&updatedResource)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedResource.ID != resourceId {
		http.Error(w, "ID mismatch between URL and request body", http.StatusBadRequest)
		return
	}

	err = RessourcesSrv.UpdateResource(resourceId, &updatedResource)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error updating resource: %v", err)
		http.Error(w, "Failed to update resource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


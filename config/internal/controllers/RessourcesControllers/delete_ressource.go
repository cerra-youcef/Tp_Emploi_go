package RessourcesControllers

import (

	"github.com/google/uuid"
	"database/sql"


	"github.com/sirupsen/logrus"

	"config/internal/services/RessourcesSrv"
	"net/http"
)



func DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la ressource depuis le contexte
	resourceIdRaw := r.Context().Value("resourceId")
	logrus.Infof("Resource ID from context: %v", resourceIdRaw)

	if resourceIdRaw == nil || resourceIdRaw.(string) == "" {
		http.Error(w, "Missing or invalid resourceId in context", http.StatusBadRequest)
		logrus.Errorf("Missing or invalid resourceId in context")
		return
	}

	// Convertir l'ID en UUID
	resourceIdStr := resourceIdRaw.(string)
	resourceId, err := uuid.Parse(resourceIdStr) // Utilisez Parse au lieu de FromString
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		logrus.Errorf("Invalid UUID format: %s", resourceIdStr)
		return
	}

	logrus.Infof("Deleting resource with ID: %s", resourceId)

	// Appeler le service pour supprimer la ressource
	err = RessourcesSrv.DeleteResource(resourceId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
		logrus.Errorf("Error deleting resource: %v", err)
		http.Error(w, "Failed to delete resource", http.StatusInternalServerError)
		return
	}

	// Répondre avec un statut 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
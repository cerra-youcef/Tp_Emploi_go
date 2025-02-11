package RessourcesControllers

import (

	"github.com/google/uuid"
	"database/sql"


	"github.com/sirupsen/logrus"

	"middleware/example/internal/services/RessourcesSrv"
	"net/http"
)

// DeleteResourceHandler supprime une ressource existante.
func DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	resourceId := r.Context().Value("resourceId").(uuid.UUID)

	err := RessourcesSrv.DeleteResource(resourceId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
		logrus.Errorf("Error deleting resource: %v", err)
		http.Error(w, "Failed to delete resource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

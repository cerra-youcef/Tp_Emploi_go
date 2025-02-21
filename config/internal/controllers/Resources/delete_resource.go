package Resources

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"config/internal/services/Resources"
)

func DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	// Safely retrieve the resource ID from context
	resourceIdRaw := r.Context().Value("resourceId")
	if resourceIdRaw == nil {
		http.Error(w, "Missing resourceId in context", http.StatusBadRequest)
		logrus.Error("Missing resourceId in context")
		return
	}

	// Ensure it's of type uuid.UUID
	resourceId, ok := resourceIdRaw.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid resourceId format", http.StatusBadRequest)
		logrus.Errorf("Invalid resourceId type: %T", resourceIdRaw)
		return
	}

	logrus.Infof("Deleting resource with ID: %s", resourceId)

	err := Resources.DeleteResource(resourceId)
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

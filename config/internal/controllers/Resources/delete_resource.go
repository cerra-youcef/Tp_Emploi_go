package Resources

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"config/internal/services/Resources"
)

// DeleteResourceHandler handles the deletion of a resource by its ID.
// @Summary Delete a resource by its ID
// @Description This endpoint allows you to delete a resource by providing the resource ID.
// @Tags Resources
// @Param resourceId path string true "Resource ID" example("123e4567-e89b-12d3-a456-426614174000")
// @Success 204 "Resource deleted successfully"
// @Failure 400 {string} string "Invalid resource ID"
// @Failure 404 {string} string "Resource not found"
// @Failure 500 {string} string "Internal server error"
// @Router /resources/{resourceId} [delete]
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

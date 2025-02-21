package Resources

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"config/internal/models"
	"config/internal/services/Resources"
	"net/http"
	"github.com/google/uuid"

)

// UpdateResourceHandler updates an existing resource.
// @Summary Update an existing resource
// @Description This endpoint allows the updating of a resource based on its ID.
// @Tags Resources
// @Accept json
// @Produce json
// @Param resourceId path string true "Resource ID"  // Param for resource ID passed in URL path
// @Param resource body models.Resource true "Updated resource object" // Resource object in the body of the request
// @Success 204 "Resource successfully updated"  // No content when update is successful
// @Failure 400 {string} string "Invalid request body or invalid resource ID"
// @Failure 404 {string} string "Resource not found"
// @Failure 500 {string} string "Failed to update resource"
// @Router /resources/{resourceId} [put]
func UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
	resourceId := r.Context().Value("resourceId").(uuid.UUID)

	var updatedResource models.Resource
	err := json.NewDecoder(r.Body).Decode(&updatedResource)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = Resources.UpdateResource(resourceId, &updatedResource)
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


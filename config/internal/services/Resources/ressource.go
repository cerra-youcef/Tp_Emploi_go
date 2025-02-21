package Resources

import (
"github.com/google/uuid"
"config/internal/models"
"config/internal/repositories/Resources"
"net/http" // Ajoutez cet import si nécessaire

)

// GetAllResources récupère toutes les ressources.
func GetAllResources() ([]models.Resource, error) {
	return Resources.GetAllResources()
}

// GetResourceByID récupère une ressource par son ID.
func GetResourceByID(id uuid.UUID) (*models.Resource, error) {
	return Resources.GetResourceByID(id)
}

// CreateResource crée une nouvelle ressource.
func CreateResource(resource *models.Resource) error {
	if resource.ID == uuid.Nil {
		resource.ID = uuid.New()
	}
	return Resources.CreateResource(resource)
}

// UpdateResource met à jour une ressource existante.
func UpdateResource(id uuid.UUID, resource *models.Resource) error {
	_, err := GetResourceByID(id)
	if err != nil {
		return &models.CustomError{
			Message: "Resource not found ",
			Code:    http.StatusNotFound,
		}
	}

	return Resources.UpdateResource(id, resource)
}

// DeleteResource supprime une ressource existante.
func DeleteResource(id uuid.UUID) error {
	return Resources.DeleteResource(id)
}

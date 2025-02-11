package RessourcesSrv

import (
"github.com/google/uuid"
"middleware/example/internal/models"
"middleware/example/internal/repositories/Ressourcesrep"
"net/http" // Ajoutez cet import si nécessaire

)

// GetAllResources récupère toutes les ressources.
func GetAllResources() ([]models.Ressource, error) {
	return Ressourcesrep.GetAllResources()
}

// GetResourceByID récupère une ressource par son ID.
func GetResourceByID(id uuid.UUID) (*models.Ressource, error) {
	return Ressourcesrep.GetResourceByID(id)
}

// CreateResource crée une nouvelle ressource.
func CreateResource(ressource *models.Ressource) error {
	if ressource.ID == uuid.Nil {
		ressource.ID = uuid.New()
	}
	return Ressourcesrep.CreateResource(ressource)
}

// UpdateResource met à jour une ressource existante.
func UpdateResource(id uuid.UUID, ressource *models.Ressource) error {
	if ressource.ID != id {
		return &models.CustomError{
			Message: "ID mismatch between URL and request body",
			Code:    http.StatusBadRequest,
		}
	}
	return Ressourcesrep.UpdateResource(id, ressource)
}

// DeleteResource supprime une ressource existante.
func DeleteResource(id uuid.UUID) error {
	return Ressourcesrep.DeleteResource(id)
}

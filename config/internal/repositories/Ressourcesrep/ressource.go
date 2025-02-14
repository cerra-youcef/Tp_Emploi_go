package Ressourcesrep

import (
	"database/sql"
	"github.com/google/uuid"
	"cerra/tp_go/internal/helpers"
	"cerra/tp_go/internal/models"
)

// GetAllResources récupère toutes les ressources.
func GetAllResources() ([]models.Ressource, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, uca_id, name
        FROM collections
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []models.Ressource
	for rows.Next() {
		var resource models.Ressource
		err = rows.Scan(&resource.ID, &resource.UcaId, &resource.Name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// GetResourceByID récupère une ressource par son ID.
func GetResourceByID(id uuid.UUID) (*models.Ressource, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, uca_id, name
        FROM collections
        WHERE id = ?
    `

	var resource models.Ressource
	err = db.QueryRow(query, id).Scan(&resource.ID, &resource.UcaId, &resource.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retourne nil si la ressource n'existe pas.
		}
		return nil, err
	}

	return &resource, nil
}

// CreateResource crée une nouvelle ressource.
func CreateResource(resource *models.Ressource) error {
	if resource.ID == uuid.Nil {
		resource.ID = uuid.New()
	}

	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        INSERT INTO collections (id, uca_id, name)
        VALUES (?, ?, ?)
    `

	result, err := db.Exec(query, resource.ID, resource.UcaId, resource.Name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// UpdateResource met à jour une ressource existante.
func UpdateResource(id uuid.UUID, resource *models.Ressource) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        UPDATE collections
        SET uca_id = ?, name = ?
        WHERE id = ?
    `

	result, err := db.Exec(query, resource.UcaId, resource.Name, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// DeleteResource supprime une ressource existante.
func DeleteResource(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        DELETE FROM collections
        WHERE id = ?
    `

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}
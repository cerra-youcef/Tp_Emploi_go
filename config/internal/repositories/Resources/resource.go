package Resources

import (
	"database/sql"
	"github.com/google/uuid"
	"config/internal/helpers"
	"config/internal/models"
)

// GetAllResources récupère toutes les ressources.
func GetAllResources() ([]models.Resource, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, uca_id, name
        FROM resources
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []models.Resource
	for rows.Next() {
		var resource models.Resource
		err = rows.Scan(&resource.ID, &resource.UcaId, &resource.Name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// GetResourceByID récupère une ressource par son ID.
func GetResourceByID(id uuid.UUID) (*models.Resource, error) {
	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `
        SELECT id, uca_id, name
        FROM resources
        WHERE id = ?
    `

	var resource models.Resource
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
func CreateResource(resource *models.Resource) error {
	if resource.ID == uuid.Nil {
		resource.ID = uuid.New()
	}

	db, err := helpers.OpenDB() // Déclaration initiale.
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        INSERT INTO resources (id, uca_id, name)
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
func UpdateResource(id uuid.UUID, resource *models.Resource) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := `
        UPDATE resources
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
        DELETE FROM resources
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
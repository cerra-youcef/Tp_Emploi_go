package models

import (
	"github.com/google/uuid"
)

type Event struct {
	ID         uuid.UUID   `json:"id"`          // Identifiant unique généré localement
	Resources []int `json:"resources"` // ID des ressource associée
	UID        string      `json:"uid"`         // UID de l'événement
	Description string     `json:"description"` // Description de l'événement
	Name       string      `json:"name"`        // Nom de l'événement
	Start      string   `json:"start"`       // Heure de début
	End        string      `json:"end"`         // Heure de fin (chaîne pour compatibilité)
	Location   string      `json:"location"`    // Lieu de l'événement
	UpdatedAt  string   `json:"updated_at"`  // Dernière mise à jour
}
package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID         uuid.UUID   `json:"id"`          // Identifiant unique généré localement
	RessourceIDs []uuid.UUID `json:"ressourceIds"` // IDs des ressources associées
	UID        string      `json:"uid"`         // UID de l'événement
	Description string     `json:"description"` // Description de l'événement
	Name       string      `json:"name"`        // Nom de l'événement
	Start      time.Time   `json:"start"`       // Heure de début
	End        string      `json:"end"`         // Heure de fin (chaîne pour compatibilité)
	Location   string      `json:"location"`    // Lieu de l'événement
	CreatedAt  time.Time   `json:"created_at"`  // Date de création
	UpdatedAt  time.Time   `json:"updated_at"`  // Dernière mise à jour
	DTStamp    string      `json:"dtStamp"`     // Horodatage
}
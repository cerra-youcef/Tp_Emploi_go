package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID           uuid.UUID   `json:"id"`
	RessourceIDs []uuid.UUID `json:"ressourceIds"`
	UID          string      `json:"uid"`
	Description  string      `json:"description"`
	Name         string      `json:"name"`
	Start        time.Time   `json:"start"`
	End         string    `json:"end"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`   // Date de création
	UpdatedAt   time.Time `json:"updated_at"`   // Dernière mise à jour
}

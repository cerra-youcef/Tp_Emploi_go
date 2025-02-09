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
}

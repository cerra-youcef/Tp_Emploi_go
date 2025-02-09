package models

import (
	"github.com/google/uuid"
)

type Alert struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	All         bool       `json:"all"`
	RessourceID *uuid.UUID `json:"ressourceId,omitempty"`
}

package models

import (
	"github.com/google/uuid"
)

type Alert struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	ResourceID string    `json:"resourceId,omitempty"`
}

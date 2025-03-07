package models

import (
	"github.com/google/uuid"
)

type Alert struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	All         bool       `json:"is_all"`
	ResourceID *uuid.UUID `json:"resourceId,omitempty"`
}

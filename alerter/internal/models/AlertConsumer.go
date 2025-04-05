package models

import (
	"github.com/google/uuid"
)

type AlertConsumer struct {
	ID      uuid.UUID         `json:"id"`
	Type    string            `json:"type"`
	Changes map[string]string `json:"changes`
	Event   uuid.UUID         `json:"event"`
}

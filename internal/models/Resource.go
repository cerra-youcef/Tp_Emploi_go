package models

import (
	"github.com/google/uuid"
)

type Ressource struct {
	ID    uuid.UUID `json:"id"`
	UcaId int       `json:"ucaId"`
	Name  string    `json:"name"`
}

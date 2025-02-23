package models

import (
	"github.com/google/uuid"
)

type Alert struct {
	ID          uuid.UUID	`json:"id"`
	Type		string		`json:"type"`
	Changes 	map[string]string		`json:"changes`	
	Event 		uuid.UUID `json:"event"`
}

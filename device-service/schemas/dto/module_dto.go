package dto

import (
	"github.com/google/uuid"
	"time"
)

type ModuleDto struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	State       string    `json:"state"`
}

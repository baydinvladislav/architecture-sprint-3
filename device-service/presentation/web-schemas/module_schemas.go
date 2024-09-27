package web_schemas

import (
	"github.com/google/uuid"
	"time"
)

type ModuleOut struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

type ConnectModuleIn struct {
	ID uuid.UUID `json:"id"`
}

type ConnectModuleOut struct {
	Modules []ModuleOut `json:"modules"`
}

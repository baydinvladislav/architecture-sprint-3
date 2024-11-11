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
	State       string    `json:"state"`
}

type ConnectModuleIn struct {
	ID uuid.UUID `json:"id"`
}

type ConnectModuleOut struct {
	Modules []ModuleOut `json:"modules"`
}

type HouseModuleState struct {
	CreatedAt time.Time `json:"created_at"`
	HouseID   uuid.UUID `json:"house_id"`
	ModuleID  uuid.UUID `json:"module_id"`
	State     string    `json:"state"`
}

package dto

import (
	"github.com/google/uuid"
)

type HouseModuleStateDto struct {
	ID       uuid.UUID         `json:"id"`
	HouseID  uuid.UUID         `json:"house_id"`
	ModuleID uuid.UUID         `json:"module_id"`
	State    map[string]string `json:"state"`
}

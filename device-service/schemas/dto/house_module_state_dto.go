package dto

import (
	"github.com/google/uuid"
	"time"
)

type HouseModuleStateDto struct {
	ID        uuid.UUID              `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	HouseID   uuid.UUID              `json:"house_id"`
	ModuleID  uuid.UUID              `json:"module_id"`
	State     map[string]interface{} `json:"state"`
}

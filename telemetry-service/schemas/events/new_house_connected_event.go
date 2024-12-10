package events

import "github.com/google/uuid"

// InstallModuleToHousePayload фиксируем ивентом иницилизации только что подключенный дом к системе SmartHome,
// в основном для реализации SAGA
type InstallModuleToHousePayload struct {
	HouseID  uint      `json:"house_id"`
	ModuleID uuid.UUID `json:"module_id"`
}

func (InstallModuleToHousePayload) IsEventPayload() {}

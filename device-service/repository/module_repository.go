package repository

import (
	"device-service/schemas/web"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrModuleAlreadyOff        = fmt.Errorf("module is already turned off")
	ErrModuleAlreadyOn         = fmt.Errorf("module is already turned on")
	ErrModuleNotFound          = fmt.Errorf("module not found")
	ErrConnectedModuleNotFound = fmt.Errorf("module not found")
)

type ModuleRepository interface {
	GetAllModules() ([]web.ModuleOut, error)
	GetModulesByHouseID(houseID uuid.UUID) ([]web.ModuleOut, error)
	TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error
	TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error
	RequestAddingModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) ([]web.ModuleOut, error)
	AcceptAdditionModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) error
	FailAdditionModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) error
	GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*web.HouseModuleState, error)
	InsertNewHouseModuleState(houseModuleId uuid.UUID, state map[string]interface{}) error
}

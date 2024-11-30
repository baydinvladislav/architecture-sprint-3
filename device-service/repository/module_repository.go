package repository

import (
	"device-service/schemas/dto"
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
	GetAllModules() ([]dto.ModuleDto, error)
	GetModulesByHouseID(houseID uuid.UUID) ([]dto.ModuleDto, error)
	TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error
	TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error
	SetPendingNewModule(houseID uuid.UUID, moduleID uuid.UUID) ([]dto.ModuleDto, error)
	AcceptAdditionModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) error
	FailAdditionModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) error
	GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*dto.HouseModuleStateDto, error)
	InsertNewHouseModuleState(houseModuleId uuid.UUID, state map[string]interface{}) error
}

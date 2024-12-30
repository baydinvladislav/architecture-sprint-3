package repository

import (
	"device-service/schemas/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockModuleRepository struct {
	mock.Mock
}

func (m *MockModuleRepository) GetAllModules() ([]dto.ModuleDto, error) {
	args := m.Called()
	return args.Get(0).([]dto.ModuleDto), args.Error(1)
}

func (m *MockModuleRepository) GetModulesByHouseID(houseID uuid.UUID) ([]dto.ModuleDto, error) {
	args := m.Called(houseID)
	return args.Get(0).([]dto.ModuleDto), args.Error(1)
}

func (m *MockModuleRepository) TurnOnModule(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) TurnOffModule(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) GetModuleState(houseID, moduleID uuid.UUID) (*dto.HouseModuleStateDto, error) {
	args := m.Called(houseID, moduleID)
	return args.Get(0).(*dto.HouseModuleStateDto), args.Error(1)
}

func (m *MockModuleRepository) AcceptAdditionModuleToHouse(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) FailAdditionModuleToHouse(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) SetPendingNewModule(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]dto.ModuleDto, error) {
	args := m.Called(houseID, moduleID)
	return args.Get(0).([]dto.ModuleDto), args.Error(1)
}

func (m *MockModuleRepository) InsertNewHouseModuleState(houseModuleId uuid.UUID, state map[string]interface{}) error {
	args := m.Called(houseModuleId, state)
	return args.Error(0)
}

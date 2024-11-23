package service

import (
	"device-service/repository"
	"device-service/schemas/dto"
	"github.com/google/uuid"
)

type ModulePersistenceService struct {
	repository repository.ModuleRepository
}

func NewModulePersistenceService(repo repository.ModuleRepository) *ModulePersistenceService {
	return &ModulePersistenceService{
		repository: repo,
	}
}

func (s *ModulePersistenceService) GetAllModules() ([]dto.ModuleDto, error) {
	return s.repository.GetAllModules()
}

func (s *ModulePersistenceService) GetModulesByHouseID(houseID uuid.UUID) ([]dto.ModuleDto, error) {
	return s.repository.GetModulesByHouseID(houseID)
}

func (s *ModulePersistenceService) GetModuleState(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) (*dto.HouseModuleStateDto, error) {
	return s.repository.GetModuleState(houseID, moduleID)
}

func (s *ModulePersistenceService) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repository.TurnOnModule(houseID, moduleID)
}

func (s *ModulePersistenceService) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repository.TurnOffModule(houseID, moduleID)
}

func (s *ModulePersistenceService) AcceptAdditionModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repository.AcceptAdditionModuleToHouse(houseID, moduleID)
}

func (s *ModulePersistenceService) FailAdditionModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repository.FailAdditionModuleToHouse(houseID, moduleID)
}

func (s *ModulePersistenceService) RequestAddingModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]dto.ModuleDto, error) {
	return s.repository.RequestAddingModuleToHouse(houseID, moduleID)
}

func (s *ModulePersistenceService) InsertNewHouseModuleState(
	houseModuleID uuid.UUID,
	state map[string]interface{},
) error {
	return s.repository.InsertNewHouseModuleState(houseModuleID, state)
}

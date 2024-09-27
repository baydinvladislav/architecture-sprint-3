package service

import (
	"device-service/presentation/web-schemas"
	"device-service/repository"
	"github.com/google/uuid"
)

type ModuleService struct {
	repo repository.ModuleRepository
}

func NewModuleService(repo repository.ModuleRepository) *ModuleService {
	return &ModuleService{
		repo: repo,
	}
}

func (s *ModuleService) GetAllModules() ([]web_schemas.ModuleOut, error) {
	return s.repo.GetAllModules()
}

func (s *ModuleService) GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	return s.repo.GetModulesByHouseID(houseID)
}

func (s *ModuleService) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repo.TurnOnModule(houseID, moduleID)
}

func (s *ModuleService) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repo.TurnOffModule(houseID, moduleID)
}

func (s *ModuleService) AddModuleToHouse(
	houseID uuid.UUID,
	newModule web_schemas.ConnectModuleIn,
) ([]web_schemas.ModuleOut, error) {
	return s.repo.AddModuleToHouse(houseID, newModule)
}

package service

import (
	"device-service/presentation/web-schemas"
	"device-service/repository"
)

type ModuleService struct {
	repo repository.GORMModuleRepository
}

func NewModuleService(repo repository.GORMModuleRepository) *ModuleService {
	return &ModuleService{
		repo: repo,
	}
}

func (s *ModuleService) GetAllModules() ([]web_schemas.ModuleOut, error) {
	return s.repo.GetAllModules()
}

func (s *ModuleService) GetModulesByHouseID(houseID uint) ([]web_schemas.ModuleOut, error) {
	return s.repo.GetModulesByHouseID(houseID)
}

func (s *ModuleService) TurnOnModule(houseID uint, moduleID uint) error {
	return s.repo.TurnOnModule(houseID, moduleID)
}

func (s *ModuleService) TurnOffModule(houseID uint, moduleID uint) error {
	return s.repo.TurnOffModule(houseID, moduleID)
}

func (s *ModuleService) AddModuleToHouse(houseID uint, newModule web_schemas.ConnectModuleIn) (*web_schemas.ModuleOut, error) {
	return s.repo.AddModuleToHouse(houseID, newModule)
}

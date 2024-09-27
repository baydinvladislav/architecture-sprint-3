package service

import (
	"device-service/presentation/web-schemas"
	"device-service/repository"
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

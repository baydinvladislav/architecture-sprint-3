package repository

import (
	"device-service/persistance"
	"device-service/presentation/web-schemas"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	GetAllModules() ([]web_schemas.ModuleOut, error)
}

type GORMModuleRepository struct {
	db *gorm.DB
}

func NewGORMModuleRepository(db *gorm.DB) *GORMModuleRepository {
	return &GORMModuleRepository{
		db: db,
	}
}

func (r *GORMModuleRepository) GetAllModules() ([]web_schemas.ModuleOut, error) {
	var modules []persistance.ModuleModel
	if err := r.db.Find(&modules).Error; err != nil {
		return nil, err
	}

	if len(modules) == 0 {
		return []web_schemas.ModuleOut{}, nil
	}

	var moduleOuts []web_schemas.ModuleOut
	for _, module := range modules {
		moduleOuts = append(moduleOuts, web_schemas.ModuleOut{
			ID:          module.ID,
			CreatedAt:   module.CreatedAt,
			Type:        module.Type,
			Description: module.Description,
		})
	}

	return moduleOuts, nil
}

package repository

import (
	"device-service/persistance"
	"device-service/presentation/web-schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	GetAllModules() ([]web_schemas.ModuleOut, error)
	GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error)
	TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error
	TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error
	AddModuleToHouse(houseID uuid.UUID, newModule web_schemas.ConnectModuleIn) ([]web_schemas.ModuleOut, error)
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

func (r *GORMModuleRepository) GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	var houseModules []persistance.HouseModuleModel
	if err := r.db.Where("house_id = ?", houseID).Find(&houseModules).Error; err != nil {
		return nil, err
	}

	var moduleOuts []web_schemas.ModuleOut
	for _, houseModule := range houseModules {
		var module persistance.ModuleModel
		if err := r.db.First(&module, "id = ?", houseModule.ModuleID).Error; err == nil {
			moduleOuts = append(moduleOuts, web_schemas.ModuleOut{
				ID:          module.ID,
				CreatedAt:   module.CreatedAt,
				Type:        module.Type,
				Description: module.Description,
			})
		}
	}

	return moduleOuts, nil
}

func (r *GORMModuleRepository) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	var houseModule persistance.HouseModuleModel
	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&houseModule).Error; err != nil {
		return err
	}

	houseModule.TurnOn = true
	return r.db.Save(&houseModule).Error
}

func (r *GORMModuleRepository) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	var houseModule persistance.HouseModuleModel
	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&houseModule).Error; err != nil {
		return err
	}

	houseModule.TurnOn = false
	return r.db.Save(&houseModule).Error
}

func (r *GORMModuleRepository) AddModuleToHouse(
	houseID uuid.UUID,
	newModule web_schemas.ConnectModuleIn,
) ([]web_schemas.ModuleOut, error) {
	module := persistance.HouseModuleModel{
		HouseID:  houseID,
		ModuleID: newModule.ID,
		TurnOn:   true,
	}

	if err := r.db.Create(&module).Error; err != nil {
		return nil, err
	}

	houseModule := persistance.HouseModuleModel{
		HouseID:  module.HouseID,
		ModuleID: module.ModuleID,
		TurnOn:   true,
	}

	if err := r.db.Create(&houseModule).Error; err != nil {
		return nil, err
	}

	modules, err := r.GetModulesByHouseID(houseID)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

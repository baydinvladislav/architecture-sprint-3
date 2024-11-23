package repository

import (
	"device-service/persistance"
	"device-service/schemas/web"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GORMModuleRepository struct {
	db *gorm.DB
}

func NewGORMModuleRepository(db *gorm.DB) *GORMModuleRepository {
	return &GORMModuleRepository{
		db: db,
	}
}

func (r *GORMModuleRepository) GetAllModules() ([]web.ModuleOut, error) {
	var modules []persistance.ModuleModel
	if err := r.db.Find(&modules).Error; err != nil {
		return nil, err
	}

	var moduleOuts []web.ModuleOut
	for _, module := range modules {
		moduleOuts = append(moduleOuts, web.ModuleOut{
			ID:          module.ID,
			CreatedAt:   module.CreatedAt,
			Type:        module.Type,
			Description: module.Description,
		})
	}

	return moduleOuts, nil
}

func (r *GORMModuleRepository) GetModulesByHouseID(houseID uuid.UUID) ([]web.ModuleOut, error) {
	var houseModules []persistance.HouseModuleModel
	if err := r.db.Where("house_id = ?", houseID).Find(&houseModules).Error; err != nil {
		return nil, err
	}

	var moduleOuts []web.ModuleOut
	for _, houseModule := range houseModules {
		var module persistance.ModuleModel
		if err := r.db.First(&module, "id = ?", houseModule.ModuleID).Error; err == nil {
			state := "activated"
			if !houseModule.TurnOn {
				state = "disabled"
			}

			moduleOuts = append(moduleOuts, web.ModuleOut{
				ID:          module.ID,
				CreatedAt:   module.CreatedAt,
				Type:        module.Type,
				Description: module.Description,
				State:       state,
			})
		}
	}

	return moduleOuts, nil
}

func (r *GORMModuleRepository) RequestAddingModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]web.ModuleOut, error) {
	var existingModule persistance.HouseModuleModel
	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&existingModule).Error; err == nil {
		return nil, fmt.Errorf("module with houseID %s and moduleID %s already exists", houseID, moduleID)
	}

	module := persistance.HouseModuleModel{
		HouseID:  houseID,
		ModuleID: moduleID,
		Status:   persistance.InstallRequested,
		TurnOn:   true,
	}

	if err := r.db.Create(&module).Error; err != nil {
		return nil, err
	}

	modules, err := r.GetModulesByHouseID(houseID)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *GORMModuleRepository) AcceptAdditionModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) error {
	var existingModule persistance.HouseModuleModel

	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&existingModule).Error; err != nil {
		return fmt.Errorf("module with houseID %s and moduleID %s not found", houseID, moduleID)
	}

	existingModule.Status = persistance.InstallCompleted
	if err := r.db.Save(&existingModule).Error; err != nil {
		return err
	}

	return nil
}

func (r *GORMModuleRepository) FailAdditionModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) error {
	var existingModule persistance.HouseModuleModel

	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&existingModule).Error; err != nil {
		return fmt.Errorf("module with houseID %s and moduleID %s not found", houseID, moduleID)
	}

	existingModule.Status = persistance.InstallFailed
	if err := r.db.Save(&existingModule).Error; err != nil {
		return err
	}

	return nil
}

func (r *GORMModuleRepository) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	var houseModule persistance.HouseModuleModel
	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&houseModule).Error; err != nil {
		return err
	}

	if houseModule.TurnOn {
		return ErrModuleAlreadyOn
	}

	houseModule.TurnOn = true
	return r.db.Save(&houseModule).Error
}

func (r *GORMModuleRepository) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	var houseModule persistance.HouseModuleModel
	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&houseModule).Error; err != nil {
		return err
	}

	if !houseModule.TurnOn {
		return ErrModuleAlreadyOff
	}

	houseModule.TurnOn = false
	return r.db.Save(&houseModule).Error
}

func (r *GORMModuleRepository) GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*web.HouseModuleState, error) {
	var houseModule persistance.HouseModuleModel

	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&houseModule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConnectedModuleNotFound
		}
		return nil, err
	}

	state := "disabled"
	if houseModule.TurnOn {
		state = "activated"
	}

	response := &web.HouseModuleState{
		ID:       houseModule.ID,
		HouseID:  houseModule.HouseID,
		ModuleID: houseModule.ModuleID,
		State:    state,
	}

	return response, nil
}

func (r *GORMModuleRepository) InsertNewHouseModuleState(houseModuleId uuid.UUID, state map[string]interface{}) error {
	houseModuleHistoryState := persistance.HouseModuleHistoryStateModel{
		ID:            uuid.New(),
		HouseModuleID: houseModuleId,
		State:         state,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := r.db.Create(&houseModuleHistoryState).Error; err != nil {
		return fmt.Errorf("failed to insert new house module state: %w", err)
	}

	return nil
}

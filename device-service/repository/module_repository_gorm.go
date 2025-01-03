package repository

import (
	"device-service/persistance"
	"device-service/schemas/dto"
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

func (r *GORMModuleRepository) GetAllModules() ([]dto.ModuleDto, error) {
	var modules []persistance.ModuleModel
	if err := r.db.Find(&modules).Error; err != nil {
		return nil, err
	}

	var modulesDto []dto.ModuleDto
	for _, module := range modules {
		modulesDto = append(modulesDto, dto.ModuleDto{
			ID:          module.ID,
			CreatedAt:   module.CreatedAt,
			Type:        module.Type,
			Description: module.Description,
		})
	}
	return modulesDto, nil
}

func (r *GORMModuleRepository) GetModulesByHouseID(houseID uuid.UUID) ([]dto.ModuleDto, error) {
	var houseModules []persistance.HouseModuleModel
	if err := r.db.Where("house_id = ?", houseID).Find(&houseModules).Error; err != nil {
		return nil, err
	}

	var modulesDto []dto.ModuleDto
	for _, houseModule := range houseModules {
		var module persistance.ModuleModel
		if err := r.db.First(&module, "id = ?", houseModule.ModuleID).Error; err == nil {
			state := "activated"
			if !houseModule.TurnOn {
				state = "disabled"
			}

			modulesDto = append(modulesDto, dto.ModuleDto{
				ID:          module.ID,
				CreatedAt:   module.CreatedAt,
				Type:        module.Type,
				Description: module.Description,
				State:       state,
			})
		}
	}
	return modulesDto, nil
}

func (r *GORMModuleRepository) SetPendingNewModule(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]dto.ModuleDto, error) {
	var existingModule persistance.HouseModuleModel
	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&existingModule).Error; err == nil {
		return nil, fmt.Errorf("module with houseID %s and moduleID %s already exists", houseID, moduleID)
	}

	module := persistance.HouseModuleModel{
		HouseID:  houseID,
		ModuleID: moduleID,
		Status:   persistance.InstallPending,
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

func (r *GORMModuleRepository) GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*dto.HouseModuleStateDto, error) {
	var houseModule persistance.HouseModuleModel

	if err := r.db.Where("house_id = ? AND module_id = ?", houseID, moduleID).First(&houseModule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConnectedModuleNotFound
		}
		return nil, err
	}

	state := map[string]interface{}{"state": "disabled"}
	if houseModule.TurnOn {
		state = map[string]interface{}{"state": "activated"}
	}

	response := &dto.HouseModuleStateDto{
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

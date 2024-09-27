package persistance

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModuleModel struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Type        string    `gorm:"size:50"`
	Description string    `gorm:"type:text"`
}

func (ModuleModel) TableName() string {
	return "modules"
}

type HouseModuleModel struct {
	HouseID  uuid.UUID `gorm:"type:uuid;not null"`
	ModuleID uuid.UUID `gorm:"type:uuid;not null"`
	TurnOn   bool
}

func (HouseModuleModel) TableName() string {
	return "house_modules"
}

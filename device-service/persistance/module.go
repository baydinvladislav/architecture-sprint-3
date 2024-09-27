package persistance

import (
	"gorm.io/gorm"
)

type ModuleModel struct {
	gorm.Model
	Type        string `gorm:"size:50"`
	Description string `gorm:"type:text"`
}

func (ModuleModel) TableName() string {
	return "modules"
}

type HouseModuleModel struct {
	HouseID  uint `gorm:"type:uuid;not null"`
	ModuleID uint `gorm:"type:uuid;not null"`
	TurnOn   bool
}

func (HouseModuleModel) TableName() string {
	return "house_modules"
}

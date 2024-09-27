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

type HouseModule struct {
	HouseID  string `gorm:"type:uuid;not null"`
	ModuleID string `gorm:"type:uuid;not null"`
	TurnOn   bool
}

func (HouseModule) TableName() string {
	return "house_modules"
}

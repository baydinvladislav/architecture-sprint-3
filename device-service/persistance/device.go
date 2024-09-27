package persistance

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Device struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	gorm.Model
	Name        string `gorm:"size:255"`
	VendorName  string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	ModuleID    string `gorm:"type:uuid;not null"`
}

func (Device) TableName() string {
	return "devices"
}

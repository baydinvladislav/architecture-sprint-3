package persistance

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `gorm:"size:255"`
	VendorName  string    `gorm:"size:255"`
	Description string    `gorm:"type:text"`
	ModuleID    uuid.UUID `gorm:"type:uuid;not null"`
}

func (Device) TableName() string {
	return "devices"
}

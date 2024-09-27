package persistance

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name        string `gorm:"size:255"`
	VendorName  string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	ModuleID    string `gorm:"type:uuid;not null"`
}

func (Device) TableName() string {
	return "devices"
}

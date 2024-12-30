package persistance

import (
	"github.com/google/uuid"
	"time"
)

type DeviceModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	Name        string     `gorm:"size:255"`
	VendorName  string     `gorm:"size:255"`
	Description string     `gorm:"type:text"`
	ModuleID    uuid.UUID  `gorm:"type:uuid;not null"`
}

func (DeviceModel) TableName() string {
	return "devices"
}

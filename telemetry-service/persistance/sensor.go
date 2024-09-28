package persistance

import (
	"github.com/google/uuid"
	"time"
)

type Sensor struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
	Type       string     `gorm:"size:50"`
	VendorName string     `gorm:"size:255"`
}

func (Sensor) TableName() string {
	return "sensors"
}

type SensorModule struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	ModuleID  uuid.UUID  `gorm:"type:uuid;not null"`
	SensorID  uuid.UUID  `gorm:"type:uuid;not null"`
	TurnOn    bool
}

func (SensorModule) TableName() string {
	return "sensor_modules"
}

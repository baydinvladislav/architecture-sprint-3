package persistance

import (
	"github.com/google/uuid"
	"time"
)

type ModuleModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	Type        string     `gorm:"size:50"`
	Description string     `gorm:"type:text"`
}

func (ModuleModel) TableName() string {
	return "modules"
}

type HouseModuleModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	HouseID   uuid.UUID  `gorm:"type:uuid;not null"`
	ModuleID  uuid.UUID  `gorm:"type:uuid;not null"`
	TurnOn    bool
}

func (HouseModuleModel) TableName() string {
	return "house_modules"
}

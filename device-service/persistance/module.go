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

type StatusEnum string

const (
	InstallRequested StatusEnum = "INSTALL_REQUESTED"
	InstallCompleted StatusEnum = "INSTALL_COMPLETED"
	InstallFailed    StatusEnum = "INSTALL_FAILED"
	Uninstall        StatusEnum = "UNINSTALL"
)

type HouseModuleModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	HouseID   uuid.UUID  `gorm:"type:uuid;not null"`
	ModuleID  uuid.UUID  `gorm:"type:uuid;not null"`
	TurnOn    bool
	Status    StatusEnum `gorm:"type:status_enum;not null"`
}

func (HouseModuleModel) TableName() string {
	return "house_modules"
}

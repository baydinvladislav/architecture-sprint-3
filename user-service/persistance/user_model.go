package persistance

import (
	"github.com/google/uuid"
	"time"
)

type UserModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Username  string     `gorm:"uniqueIndex;size:100"`
	Password  string     `gorm:"size:255"`
}

func (UserModel) TableName() string {
	return "users"
}

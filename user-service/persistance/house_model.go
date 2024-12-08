package persistance

import (
	"github.com/google/uuid"
	"time"
)

type HouseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Address   string     `gorm:"size:255"`
	Square    float64    `gorm:"type:decimal(10,2)"`
	UserID    uuid.UUID  `gorm:"type:int;not null"`

	User UserModel `gorm:"foreignKey:UserID"`
}

func (HouseModel) TableName() string {
	return "houses"
}

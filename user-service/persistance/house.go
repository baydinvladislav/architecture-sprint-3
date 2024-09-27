package persistance

import (
	"gorm.io/gorm"
)

type HouseModel struct {
	gorm.Model
	Address string  `gorm:"size:255"`
	Square  float64 `gorm:"type:decimal(10,2)"`
	UserID  uint    `gorm:"type:int;not null"`

	User UserModel `gorm:"foreignKey:UserID"`
}

func (HouseModel) TableName() string {
	return "houses"
}

package repository

import (
	"gorm.io/gorm"
	"user-service/persistance"
	"user-service/presentation/web-schemas"
)

type HouseRepository interface {
	CreateUserHouse(house web_schemas.NewHouseIn) (*persistance.HouseModel, error)
	GetUserHouse(userID string) (web_schemas.HouseOut, error)
	UpdateUserHouse(house web_schemas.UpdateHouseIn) (*persistance.HouseModel, error)
}

type GORMHouseRepository struct {
	db *gorm.DB
}

func NewGORMHouseRepository(db *gorm.DB) *GORMHouseRepository {
	return &GORMHouseRepository{
		db: db,
	}
}

func (r *GORMHouseRepository) CreateUserHouse(house web_schemas.NewHouseIn) (*persistance.HouseModel, error) {
	newHouse := persistance.HouseModel{
		Address: house.Address,
		Square:  house.Square,
		UserID:  house.UserID,
	}

	if err := r.db.Create(&newHouse).Error; err != nil {
		return nil, err
	}

	return &newHouse, nil
}

func (r *GORMHouseRepository) GetUserHouse(userID string) (web_schemas.HouseOut, error) {
	var house persistance.HouseModel
	err := r.db.Where("user_id = ?", userID).First(&house).Error
	if err != nil {
		return web_schemas.HouseOut{}, err
	}

	return web_schemas.HouseOut{
		ID:      house.ID,
		Address: house.Address,
		Square:  house.Square,
		UserID:  house.UserID,
	}, nil
}

func (r *GORMHouseRepository) UpdateUserHouse(house web_schemas.UpdateHouseIn) (*persistance.HouseModel, error) {
	updatedHouse := persistance.HouseModel{
		Address: house.Address,
		Square:  house.Square,
	}

	var existingHouse persistance.HouseModel
	if err := r.db.First(&existingHouse, house.HouseID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingHouse).
		Where("id = ? AND user_id = ?", house.HouseID, house.UserID).
		Updates(updatedHouse).
		Error; err != nil {
		return nil, err
	}

	return &existingHouse, nil
}

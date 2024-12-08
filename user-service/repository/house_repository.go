package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"user-service/persistance"
	"user-service/schemas/web"
)

type HouseRepository interface {
	CreateUserHouse(userId uuid.UUID, house web.NewHouseIn) (*persistance.HouseModel, error)
	GetUserHouses(userID uuid.UUID) ([]web.HouseOut, error)
	UpdateUserHouse(house web.UpdateHouseIn) (*persistance.HouseModel, error)
}

type GORMHouseRepository struct {
	db *gorm.DB
}

func NewGORMHouseRepository(db *gorm.DB) *GORMHouseRepository {
	return &GORMHouseRepository{
		db: db,
	}
}

func (r *GORMHouseRepository) CreateUserHouse(
	userId uuid.UUID,
	house web.NewHouseIn,
) (*persistance.HouseModel, error) {
	newHouse := persistance.HouseModel{
		Address: house.Address,
		Square:  house.Square,
		UserID:  userId,
	}

	if err := r.db.Create(&newHouse).Error; err != nil {
		return nil, err
	}

	return &newHouse, nil
}

func (r *GORMHouseRepository) GetUserHouses(userID uuid.UUID) ([]web.HouseOut, error) {
	var houses []persistance.HouseModel
	err := r.db.Where("user_id = ?", userID).Find(&houses).Error
	if err != nil {
		return nil, err
	}

	var houseOuts []web.HouseOut
	for _, house := range houses {
		houseOuts = append(houseOuts, web.HouseOut{
			ID:      house.ID,
			Address: house.Address,
			Square:  house.Square,
			UserID:  house.UserID,
		})
	}

	return houseOuts, nil
}

func (r *GORMHouseRepository) UpdateUserHouse(house web.UpdateHouseIn) (*persistance.HouseModel, error) {
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

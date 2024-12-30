package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"user-service/persistance"
	"user-service/schemas/dto"
	"user-service/schemas/web"
)

type HouseRepository interface {
	CreateUserHouse(userId uuid.UUID, house web.NewHouseIn) (*dto.HouseDtoSchema, error)
	GetUserHouses(userID uuid.UUID) ([]dto.HouseDtoSchema, error)
	UpdateUserHouse(house web.UpdateHouseIn) (*dto.HouseDtoSchema, error)
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
) (*dto.HouseDtoSchema, error) {
	newHouse := persistance.HouseModel{
		Address: house.Address,
		Square:  house.Square,
		UserID:  userId,
	}

	if err := r.db.Create(&newHouse).Error; err != nil {
		return nil, err
	}

	houseDTO := &dto.HouseDtoSchema{
		ID:      newHouse.ID,
		Address: newHouse.Address,
		Square:  newHouse.Square,
		UserID:  newHouse.UserID,
	}

	return houseDTO, nil
}

func (r *GORMHouseRepository) GetUserHouses(userID uuid.UUID) ([]dto.HouseDtoSchema, error) {
	var houses []persistance.HouseModel
	err := r.db.Where("user_id = ?", userID).Find(&houses).Error
	if err != nil {
		return nil, err
	}

	var houseDTOs []dto.HouseDtoSchema
	for _, house := range houses {
		houseDTOs = append(houseDTOs, dto.HouseDtoSchema{
			ID:      house.ID,
			Address: house.Address,
			Square:  house.Square,
			UserID:  house.UserID,
		})
	}
	return houseDTOs, nil
}

func (r *GORMHouseRepository) UpdateUserHouse(house web.UpdateHouseIn) (*dto.HouseDtoSchema, error) {
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

	houseDto := &dto.HouseDtoSchema{
		ID:      existingHouse.ID,
		Address: updatedHouse.Address,
		Square:  updatedHouse.Square,
		UserID:  existingHouse.UserID,
	}

	return houseDto, nil
}

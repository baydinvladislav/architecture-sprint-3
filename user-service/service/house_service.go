package service

import (
	"user-service/persistance"
	web_schemas "user-service/presentation/web-schemas"
	"user-service/repository"
)

type HouseService struct {
	repo repository.HouseRepository
}

func NewHouseService(repo repository.HouseRepository) *HouseService {
	return &HouseService{
		repo: repo,
	}
}

func (s *HouseService) CreateUserHouse(house web_schemas.NewHouseIn) (*web_schemas.HouseOut, error) {
	newHouse, err := s.repo.CreateUserHouse(house)
	if err != nil {
		return nil, err
	}

	houseOut := &web_schemas.HouseOut{
		ID:      newHouse.ID,
		Address: newHouse.Address,
		Square:  newHouse.Square,
		UserID:  newHouse.UserID,
	}

	return houseOut, nil
}

func (s *HouseService) GetUserHouse(userID string) (web_schemas.HouseOut, error) {
	return s.repo.GetUserHouse(userID)
}

func (s *HouseService) UpdateUserHouse(house web_schemas.UpdateHouseIn) (*persistance.HouseModel, error) {
	return s.repo.UpdateUserHouse(house)
}

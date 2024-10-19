package service

import (
	"encoding/json"
	"fmt"
	"log"
	"user-service/persistance"
	web_schemas "user-service/presentation/web-schemas"
	"user-service/repository"
	"user-service/schemas"
	"user-service/suppliers"
)

type HouseService struct {
	repo          repository.HouseRepository
	kafkaSupplier suppliers.KafkaSupplier
	userService   UserService
}

func NewHouseService(repo repository.HouseRepository) *HouseService {
	return &HouseService{
		repo: repo,
	}
}

func (s *HouseService) CreateUserHouse(
	userId uint,
	house web_schemas.NewHouseIn,
) (*web_schemas.HouseOut, error) {
	newHouse, err := s.repo.CreateUserHouse(userId, house)
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

func (s *HouseService) GetUserHouses(userID uint) ([]web_schemas.HouseOut, error) {
	return s.repo.GetUserHouses(userID)
}

func (s *HouseService) UpdateUserHouse(house web_schemas.UpdateHouseIn) (*persistance.HouseModel, error) {
	return s.repo.UpdateUserHouse(house)
}

func (s *HouseService) ProcessEvent(event schemas.Event) error {
	var data schemas.ModuleVerifyPayload

	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		return fmt.Errorf("failed to unmarshal payload to TelemetryPayload: %v", err)
	}

	log.Println("msg data: ", data)

	if err != nil {
		return err
	}
	return nil
}

func (s *HouseService) verifyUserAndHouse(userId uint, houseId uint) (bool, error) {
	user, err := s.userService.GetRequiredById(userId)
	if err != nil {
		return false, err
	}

	if user.Username == "" {
		return false, fmt.Errorf("user with ID %d does not have a username", userId)
	}

	houses, err := s.GetUserHouses(userId)
	if err != nil {
		return false, err
	}

	var foundHouse *web_schemas.HouseOut

	for _, house := range houses {
		if house.ID == houseId {
			foundHouse = &house
			break
		}
	}

	if foundHouse.Square < 100 {
		return false, nil
	}

	return true, nil
}

func (s *HouseService) ApproveModuleInstallation(userId uint, houseId uint) (bool, error) {
	return s.verifyUserAndHouse(userId, houseId)
}

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"user-service/repository"
	"user-service/schemas/events"
	"user-service/schemas/web"
	"user-service/suppliers"
)

type HouseService struct {
	houseRepository repository.HouseRepository
	userService     *UserService
	verifyService   *VerifyConnectionService
	kafkaSupplier   *suppliers.KafkaSupplier
}

func NewHouseService(
	houseRepository repository.HouseRepository,
	userService *UserService,
	verifyService *VerifyConnectionService,
	supplier *suppliers.KafkaSupplier,
) *HouseService {
	return &HouseService{
		houseRepository: houseRepository,
		verifyService:   verifyService,
		userService:     userService,
		kafkaSupplier:   supplier,
	}
}

func (s *HouseService) CreateUserHouse(
	userId uuid.UUID,
	house web.NewHouseIn,
) (*web.HouseOut, error) {
	newHouse, err := s.houseRepository.CreateUserHouse(userId, house)
	if err != nil {
		return nil, err
	}

	return &web.HouseOut{
		ID:      newHouse.ID,
		Address: newHouse.Address,
		Square:  newHouse.Square,
		UserID:  newHouse.UserID,
	}, nil
}

func (s *HouseService) GetUserHouses(userID uuid.UUID) ([]web.HouseOut, error) {
	housesDto, err := s.houseRepository.GetUserHouses(userID)
	if err != nil {
		return nil, err
	}

	var houses []web.HouseOut
	for _, house := range housesDto {
		houses = append(houses, web.HouseOut{
			ID:      house.ID,
			Address: house.Address,
			Square:  house.Square,
			UserID:  house.UserID,
		})
	}
	return houses, nil
}

func (s *HouseService) UpdateUserHouse(house web.UpdateHouseIn) (*web.HouseOut, error) {
	updatedHouse, err := s.houseRepository.UpdateUserHouse(house)
	if err != nil {
		return &web.HouseOut{}, err
	}

	return &web.HouseOut{
		ID:      updatedHouse.ID,
		Address: updatedHouse.Address,
		Square:  updatedHouse.Square,
		UserID:  updatedHouse.UserID,
	}, nil
}

func (s *HouseService) ApproveModuleInstallation(userId uuid.UUID, houseId uuid.UUID) (bool, error) {
	verifyingUser, err := s.userService.GetRequiredById(userId)
	if err != nil {
		return false, err
	}

	userHouses, err := s.GetUserHouses(userId)
	if err != nil {
		return false, err
	}

	var verifyingHouse *web.HouseOut
	for _, house := range userHouses {
		if house.ID == houseId {
			verifyingHouse = &house
			break
		}
	}

	if verifyingHouse.ID.String() == "" {
		return false, errors.New("verifying house not found")
	}

	return s.verifyService.VerifyModuleConnection(verifyingUser, verifyingHouse)
}

func (s *HouseService) GetModuleAdditionEvent(ctx context.Context) (events.BaseEvent, error) {
	msg, err := s.kafkaSupplier.ReadModuleAdditionTopic(ctx)
	if err != nil {
		return events.BaseEvent{}, fmt.Errorf("failed to read message: %w", err)
	}

	var event events.BaseEvent
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return events.BaseEvent{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return event, nil
}

func (s *HouseService) ProcessModuleAdditionEvent(event events.BaseEvent) error {
	switch event.EventType {
	case "ModuleAdditionEvent":
		payload, ok := event.Payload.(events.ModuleAdditionEvent)
		if !ok {
			return errors.New("invalid payload type")
		}

		houseID, err := uuid.Parse(payload.HouseID)
		if err != nil {
			return errors.New("invalid houseID UUID")
		}
		moduleID, err := uuid.Parse(payload.ModuleID)
		if err != nil {
			return errors.New("invalid moduleID UUID")
		}

		decision, err := s.ApproveModuleInstallation(houseID, moduleID)
		if err != nil {
			return err
		}

		if decision == true {
			fmt.Printf(
				"Verification of the house %v for module %v is success\n", houseID, moduleID,
			)
		} else {
			fmt.Printf(
				"Verification of the house %v for module %v is failed\n", houseID, moduleID,
			)
		}
	}
	return errors.New("unsupported event type")
}

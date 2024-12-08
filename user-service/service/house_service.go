package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"user-service/persistance"
	"user-service/repository"
	"user-service/schemas/events"
	web_schemas "user-service/schemas/web"
	"user-service/suppliers"
)

type HouseService struct {
	houseRepository repository.HouseRepository
	kafkaSupplier   *suppliers.KafkaSupplier
	userService     UserService
}

func NewHouseService(repository repository.HouseRepository, supplier *suppliers.KafkaSupplier) *HouseService {
	return &HouseService{
		houseRepository: repository,
		kafkaSupplier:   supplier,
	}
}

func (s *HouseService) CreateUserHouse(
	userId uuid.UUID,
	house web_schemas.NewHouseIn,
) (*web_schemas.HouseOut, error) {
	newHouse, err := s.houseRepository.CreateUserHouse(userId, house)
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

func (s *HouseService) GetUserHouses(userID uuid.UUID) ([]web_schemas.HouseOut, error) {
	return s.houseRepository.GetUserHouses(userID)
}

func (s *HouseService) UpdateUserHouse(house web_schemas.UpdateHouseIn) (*persistance.HouseModel, error) {
	return s.houseRepository.UpdateUserHouse(house)
}

func (s *HouseService) verifyUserAndHouse(userId uuid.UUID, houseId uuid.UUID) (bool, error) {
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

func (s *HouseService) ApproveModuleInstallation(userId uuid.UUID, houseId uuid.UUID) (bool, error) {
	return s.verifyUserAndHouse(userId, houseId)
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

		decision, err := s.verifyUserAndHouse(houseID, moduleID)
		if err != nil {
			return err
		}

		fmt.Printf(
			"Verification of the house %v for module %v is complete, solution: %v\n", houseID, moduleID, decision,
		)

		return nil
	}
	return errors.New("unsupported event type")
}

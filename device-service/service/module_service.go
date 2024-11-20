package service

import (
	"context"
	"device-service/presentation/web-schemas"
	"device-service/repository"
	"device-service/schemas"
	"device-service/suppliers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ModuleService struct {
	repo          repository.ModuleRepository
	kafkaSupplier suppliers.KafkaSupplierInterface
}

var ErrKafkaSupplier = fmt.Errorf("erorr during send message in kafka")

func NewModuleService(repo repository.ModuleRepository, kafkaSupplier suppliers.KafkaSupplierInterface) *ModuleService {
	return &ModuleService{
		repo:          repo,
		kafkaSupplier: kafkaSupplier,
	}
}

func (s *ModuleService) ProcessMessage(event schemas.BaseEvent) (bool, error) {
	switch event.EventType {
	case "ModuleVerificationEvent":
		payload, ok := event.Payload.(schemas.ModuleVerificationEvent)
		if !ok {
			return false, errors.New("invalid payload type")
		}

		houseID, err := uuid.Parse(payload.HouseID)
		if err != nil {
			return false, errors.New("invalid houseID UUID")
		}
		moduleID, err := uuid.Parse(payload.ModuleID)
		if err != nil {
			return false, errors.New("invalid moduleID UUID")
		}

		if payload.Decision == "ACCEPTED" {
			err = s.acceptModuleAddition(houseID, moduleID)
			if err != nil {
				return false, errors.New("failed to accept module addition")
			}
		} else if payload.Decision == "FAILED" {
			err = s.failModuleAddition(houseID, moduleID)
			if err != nil {
				return false, errors.New("failed to process module failure")
			}
		} else {
			return false, errors.New("unsupported decision type")
		}
		return true, nil

	default:
		return false, errors.New("unsupported event type")
	}
}

func (s *ModuleService) GetModuleVerificationEvent(ctx context.Context) (schemas.BaseEvent, error) {
	msg, err := s.kafkaSupplier.ReadModuleVerificationTopic(ctx)
	if err != nil {
		return schemas.BaseEvent{}, err
	}

	var event schemas.BaseEvent
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return schemas.BaseEvent{}, errors.New("failed to unmarshal message to BaseEvent")
	}

	return event, nil
}

func (s *ModuleService) GetAllModules() ([]web_schemas.ModuleOut, error) {
	return s.repo.GetAllModules()
}

func (s *ModuleService) GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	return s.repo.GetModulesByHouseID(houseID)
}

func (s *ModuleService) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	err := s.repo.TurnOnModule(houseID, moduleID)
	if err != nil {
		return err
	}

	key := []byte(moduleID.String())
	event := schemas.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State: map[string]interface{}{
			"running": "on",
		},
	}

	if err := s.kafkaSupplier.SendMessageToEquipmentChangeStateTopic(context.Background(), key, event); err != nil {
		return err
	}

	return nil
}

func (s *ModuleService) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	err := s.repo.TurnOffModule(houseID, moduleID)
	if err != nil {
		return err
	}

	key := []byte(moduleID.String())
	event := schemas.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State: map[string]interface{}{
			"running": "off",
		},
	}

	if err := s.kafkaSupplier.SendMessageToEquipmentChangeStateTopic(context.Background(), key, event); err != nil {
		return err
	}

	return nil
}

func (s *ModuleService) GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*web_schemas.HouseModuleState, error) {
	return s.repo.GetModuleState(houseID, moduleID)
}

func (s *ModuleService) RequestAdditionModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]web_schemas.ModuleOut, error) {
	response, err := s.repo.RequestAddingModuleToHouse(houseID, moduleID)

	key := []byte(moduleID.String())
	event := schemas.HomeVerificationEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
	}

	if err := s.kafkaSupplier.SendMessageToAdditionTopic(context.Background(), key, event); err != nil {
		return nil, err
	}

	return response, err
}

func (s *ModuleService) ChangeEquipmentState(
	houseID uuid.UUID,
	moduleID uuid.UUID,
	state map[string]interface{},
) (*web_schemas.HouseModuleState, error) {
	houseModule, err := s.repo.GetModuleState(houseID, moduleID)

	if err != nil {
		if errors.Is(err, repository.ErrConnectedModuleNotFound) {
			return nil, repository.ErrConnectedModuleNotFound
		}

		return nil, fmt.Errorf("failed to get module state: %w", err)
	}

	err = s.repo.InsertNewHouseModuleState(houseModule.ID, state)
	if err != nil {
		return nil, err
	}

	key := []byte(moduleID.String())
	event := schemas.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    state,
	}

	if err := s.kafkaSupplier.SendMessageToEquipmentChangeStateTopic(context.Background(), key, event); err != nil {
		return nil, err
	}

	return houseModule, nil
}

func (s *ModuleService) acceptModuleAddition(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) error {
	err := s.repo.AcceptAdditionModuleToHouse(houseID, moduleID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ModuleService) failModuleAddition(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) error {
	err := s.repo.FailAdditionModuleToHouse(houseID, moduleID)
	if err != nil {
		return err
	}
	return nil
}

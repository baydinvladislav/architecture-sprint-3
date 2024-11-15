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
	kafkaSupplier suppliers.KafkaSupplier
}

var ErrKafkaSupplier = fmt.Errorf("erorr during send message in kafka")

func NewModuleService(repo repository.ModuleRepository) *ModuleService {
	return &ModuleService{
		repo: repo,
	}
}

func (s *ModuleService) ProcessMessage(event schemas.Event) (bool, error) {
	switch event.EventType {
	case "ModuleVerification":
		payload, ok := event.Payload.(schemas.ModuleVerification)
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

func (s *ModuleService) ReadMessage(ctx context.Context) (schemas.Event, error) {
	msg, err := s.kafkaSupplier.ReadModuleVerificationTopic(ctx)
	if err != nil {
		return schemas.Event{}, err
	}

	var event schemas.Event
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return schemas.Event{}, errors.New("failed to unmarshal message to Event")
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
	event := schemas.ChangeEquipmentState{
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
	event := schemas.ChangeEquipmentState{
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
	return s.repo.RequestAddingModuleToHouse(houseID, moduleID)
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

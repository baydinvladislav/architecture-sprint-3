package service

import (
	"context"
	web_schemas "device-service/presentation/web-schemas"
	"device-service/schemas"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ModuleService1 struct {
	persistenceService *ModulePersistenceService
	messagingService   *ExternalSystemMessagingService
}

var ErrKafkaSupplier1 = fmt.Errorf("error during send message in Kafka")

func NewModuleService1(persistenceService *ModulePersistenceService, messagingService *ExternalSystemMessagingService) *ModuleService1 {
	return &ModuleService1{
		persistenceService: persistenceService,
		messagingService:   messagingService,
	}
}

func (s *ModuleService1) ProcessMessage(event schemas.BaseEvent) (bool, error) {
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
			return true, s.persistenceService.AcceptAdditionModuleToHouse(houseID, moduleID)
		} else if payload.Decision == "FAILED" {
			return true, s.persistenceService.FailAdditionModuleToHouse(houseID, moduleID)
		}

		return false, errors.New("unsupported decision type")
	}

	return false, errors.New("unsupported event type")
}

func (s *ModuleService1) GetModuleVerificationEvent(ctx context.Context) (schemas.BaseEvent, error) {
	return s.messagingService.ReadModuleVerificationEvent(ctx)
}

func (s *ModuleService1) GetAllModules() ([]web_schemas.ModuleOut, error) {
	return s.persistenceService.GetAllModules()
}

func (s *ModuleService1) GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*web_schemas.HouseModuleState, error) {
	return s.persistenceService.GetModuleState(houseID, moduleID)
}

func (s *ModuleService1) GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	return s.persistenceService.GetModulesByHouseID(houseID)
}

func (s *ModuleService1) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	err := s.persistenceService.TurnOnModule(houseID, moduleID)
	if err != nil {
		return err
	}

	state := map[string]interface{}{"running": "on"}
	event := schemas.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    state,
	}

	return s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
}

func (s *ModuleService1) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	err := s.persistenceService.TurnOffModule(houseID, moduleID)
	if err != nil {
		return err
	}

	state := map[string]interface{}{"running": "off"}
	event := schemas.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    state,
	}

	return s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
}

func (s *ModuleService1) RequestAdditionModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]web_schemas.ModuleOut, error) {
	response, err := s.persistenceService.RequestAddingModuleToHouse(houseID, moduleID)
	if err != nil {
		return nil, err
	}

	event := schemas.HomeVerificationEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
	}

	err = s.messagingService.SendModuleAdditionEvent(context.Background(), []byte(moduleID.String()), event)
	return response, err
}

func (s *ModuleService1) ChangeEquipmentState(
	houseID uuid.UUID,
	moduleID uuid.UUID,
	state map[string]interface{},
) (*web_schemas.HouseModuleState, error) {
	houseModule, err := s.persistenceService.GetModuleState(houseID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module state: %w", err)
	}

	err = s.persistenceService.InsertNewHouseModuleState(houseModule.ID, state)
	if err != nil {
		return nil, err
	}

	event := schemas.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    state,
	}

	err = s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
	return houseModule, err
}

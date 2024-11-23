package service

import (
	"context"
	"device-service/schemas/events"
	web_schemas "device-service/schemas/web"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ModuleService struct {
	persistenceService *ModulePersistenceService
	messagingService   *ExternalMessagingService
}

func NewModuleService(
	persistenceService *ModulePersistenceService,
	messagingService *ExternalMessagingService,
) *ModuleService {
	return &ModuleService{
		persistenceService: persistenceService,
		messagingService:   messagingService,
	}
}

func (s *ModuleService) ProcessMessage(event events.BaseEvent) (bool, error) {
	switch event.EventType {
	case "ModuleVerificationEvent":
		payload, ok := event.Payload.(events.ModuleVerificationEvent)
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

func (s *ModuleService) GetModuleVerificationEvent(ctx context.Context) (events.BaseEvent, error) {
	return s.messagingService.ReadModuleVerificationEvent(ctx)
}

func (s *ModuleService) GetAllModules() ([]web_schemas.ModuleOut, error) {
	modulesDto, err := s.persistenceService.GetAllModules()
	if err != nil {
		return nil, err
	}

	var modulesOut []web_schemas.ModuleOut
	for _, m := range modulesDto {
		moduleOut := web_schemas.ModuleOut{
			ID:          m.ID,
			CreatedAt:   m.CreatedAt,
			Type:        m.Type,
			Description: m.Description,
			State:       m.State,
		}
		modulesOut = append(modulesOut, moduleOut)
	}
	return modulesOut, nil
}

func (s *ModuleService) GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	return s.persistenceService.GetModulesByHouseID(houseID)
}

func (s *ModuleService) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	err := s.persistenceService.TurnOnModule(houseID, moduleID)
	if err != nil {
		return err
	}

	moduleState, err := s.persistenceService.GetModuleState(houseID, moduleID)
	if err != nil {
		return err
	}

	newState := map[string]interface{}{"running": "on"}

	err = s.persistenceService.InsertNewHouseModuleState(moduleState.ID, newState)
	if err != nil {
		return err
	}

	event := events.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    newState,
	}

	return s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
}

func (s *ModuleService) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	err := s.persistenceService.TurnOffModule(houseID, moduleID)
	if err != nil {
		return err
	}

	moduleState, err := s.persistenceService.GetModuleState(houseID, moduleID)
	if err != nil {
		return err
	}

	newState := map[string]interface{}{"running": "off"}

	err = s.persistenceService.InsertNewHouseModuleState(moduleState.ID, newState)
	if err != nil {
		return err
	}

	event := events.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    newState,
	}

	return s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
}

func (s *ModuleService) GetModuleState(houseID uuid.UUID, moduleID uuid.UUID) (*web_schemas.HouseModuleState, error) {
	return s.persistenceService.GetModuleState(houseID, moduleID)
}

func (s *ModuleService) RequestAdditionModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]web_schemas.ModuleOut, error) {
	response, err := s.persistenceService.RequestAddingModuleToHouse(houseID, moduleID)
	if err != nil {
		return nil, err
	}

	event := events.HomeVerificationEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
	}

	err = s.messagingService.SendModuleAdditionEvent(context.Background(), []byte(moduleID.String()), event)
	return response, err
}

func (s *ModuleService) ChangeEquipmentState(
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

	event := events.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    state,
	}

	err = s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
	return houseModule, err
}

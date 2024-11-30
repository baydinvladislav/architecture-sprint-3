package service

import (
	"context"
	"device-service/repository"
	"device-service/schemas/events"
	web_schemas "device-service/schemas/web"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
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

func (s *ModuleService) GetAllModules() ([]web_schemas.ModuleOut, error) {
	modulesDto, err := s.persistenceService.GetAllModules()
	if err != nil {
		return nil, err
	}

	log.Printf("Got smart home modules: %d", len(modulesDto))

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
	modulesDto, err := s.persistenceService.GetModulesByHouseID(houseID)
	if err != nil {
		return nil, err
	}

	log.Printf("Got %d installed modules for %v house", houseID, len(modulesDto))

	var modulesOut []web_schemas.ModuleOut
	for _, m := range modulesDto {
		modulesOut = append(modulesOut, web_schemas.ModuleOut{
			ID:          m.ID,
			CreatedAt:   m.CreatedAt,
			Type:        m.Type,
			Description: m.Description,
			State:       m.State,
		})
	}
	return modulesOut, nil
}

func (s *ModuleService) RequestModuleInstallation(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]web_schemas.ModuleOut, error) {
	modulesDto, err := s.persistenceService.SetPendingNewModule(houseID, moduleID)
	if err != nil {
		return nil, err
	}

	event := events.HomeVerificationEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
	}
	err = s.messagingService.SendModuleAdditionEvent(
		context.Background(),
		[]byte(moduleID.String()),
		event,
	)

	var modulesOut []web_schemas.ModuleOut
	for _, m := range modulesDto {
		modulesOut = append(modulesOut, web_schemas.ModuleOut{
			ID:          m.ID,
			CreatedAt:   m.CreatedAt,
			Type:        m.Type,
			Description: m.Description,
			State:       m.State,
		})
	}
	return modulesOut, err
}

func (s *ModuleService) GetModuleState(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) (*web_schemas.HouseModuleState, error) {
	moduleStateDto, err := s.persistenceService.GetModuleState(houseID, moduleID)
	if err != nil {
		return nil, err
	}

	log.Printf("Got installed module %v state in user house %v", houseID, moduleID)

	moduleState := &web_schemas.HouseModuleState{
		ID:       moduleStateDto.ID,
		HouseID:  moduleStateDto.HouseID,
		ModuleID: moduleStateDto.ModuleID,
		State:    moduleStateDto.State,
	}

	return moduleState, nil
}

func (s *ModuleService) ChangeModuleState(
	houseID uuid.UUID,
	moduleID uuid.UUID,
	state map[string]interface{},
) (*web_schemas.HouseModuleState, error) {
	fmt.Printf(
		"Starting change installed module state %v in house %v with state: %v",
		houseID, moduleID, state,
	)

	houseModuleStateDto, err := s.persistenceService.GetModuleState(houseID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module state: %w", err)
	}

	err = s.persistenceService.InsertNewHouseModuleState(houseModuleStateDto.ID, state)
	if err != nil {
		return nil, err
	}

	event := events.ChangeEquipmentStateEvent{
		HouseID:  houseID.String(),
		ModuleID: moduleID.String(),
		Time:     time.Now().Unix(),
		State:    state,
	}
	houseModuleOut := &web_schemas.HouseModuleState{
		ID:        houseModuleStateDto.ID,
		CreatedAt: houseModuleStateDto.CreatedAt,
		HouseID:   houseModuleStateDto.HouseID,
		ModuleID:  houseModuleStateDto.ModuleID,
		State:     houseModuleStateDto.State,
	}

	err = s.messagingService.SendEquipmentStateChangeEvent(
		context.Background(),
		[]byte(moduleID.String()),
		event,
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf(
		"Change module %v state in house %v completed successfully, the current state is %v",
		houseID, moduleID, state,
	)
	return houseModuleOut, nil
}

func (s *ModuleService) TurnOnModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	fmt.Printf("Starting activating installed module %v in house %v", houseID, moduleID)

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

	err = s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
	if err != nil {
		return err
	}

	fmt.Printf("Installed module %v in house %v successfully activated", houseID, moduleID)
	return nil
}

func (s *ModuleService) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	fmt.Printf("Starting disabling installed module %v in house %v", houseID, moduleID)

	err := s.persistenceService.TurnOffModule(houseID, moduleID)
	if err != nil {
		if errors.Is(err, repository.ErrModuleAlreadyOff) {
			fmt.Printf("Installed module %v in house %v already disabled", houseID, moduleID)
		}
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

	err = s.messagingService.SendEquipmentStateChangeEvent(context.Background(), []byte(moduleID.String()), event)
	if err != nil {
		return err
	}

	fmt.Printf("Instaled module %v in house %v successfully disabled", houseID, moduleID)
	return nil
}

func (s *ModuleService) GetModuleVerificationEvent(ctx context.Context) (events.BaseEvent, error) {
	return s.messagingService.ReadModuleVerificationEvent(ctx)
}

func (s *ModuleService) ProcessModuleVerificationEvent(event events.BaseEvent) error {
	switch event.EventType {
	case "ModuleVerificationEvent":
		payload, ok := event.Payload.(events.ModuleVerificationEvent)
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

		decisionHandlers := map[string]func(uuid.UUID, uuid.UUID) error{
			"ACCEPTED": s.persistenceService.AcceptAdditionModuleToHouse,
			"FAILED":   s.persistenceService.FailAdditionModuleToHouse,
		}

		handler, exists := decisionHandlers[payload.Decision]
		if !exists {
			return errors.New("unsupported decision type")
		}
		return handler(houseID, moduleID)
	}
	return errors.New("unsupported event type")
}

package service

import (
	"context"
	"device-service/presentation/web-schemas"
	"device-service/repository"
	"device-service/schemas"
	"device-service/suppliers"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

type ModuleService struct {
	repo          repository.ModuleRepository
	kafkaSupplier suppliers.KafkaSupplier
}

func NewModuleService(repo repository.ModuleRepository) *ModuleService {
	return &ModuleService{
		repo: repo,
	}
}

func (s *ModuleService) ProcessMessage(event schemas.Event) (bool, error) {
	switch event.EventType {
	case "AddModuleToHouse":
		payload, ok := event.Payload.(schemas.ModuleAdditionPayload)
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

		_, err = s.AddModuleToHouse(houseID, moduleID)
		if err != nil {
			return false, errors.New("failed to add module to house")
		}

		return true, nil
	default:
		return false, errors.New("unsupported event type")
	}
}

func (s *ModuleService) ReadMessage(ctx context.Context) (schemas.Event, error) {
	msg, err := s.kafkaSupplier.ReadMessage(ctx)
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
	return s.repo.TurnOnModule(houseID, moduleID)
}

func (s *ModuleService) TurnOffModule(houseID uuid.UUID, moduleID uuid.UUID) error {
	return s.repo.TurnOffModule(houseID, moduleID)
}

func (s *ModuleService) AddModuleToHouse(
	houseID uuid.UUID,
	moduleID uuid.UUID,
) ([]web_schemas.ModuleOut, error) {
	return s.repo.AddModuleToHouse(houseID, moduleID)
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"telemetry-service/repository"
	"telemetry-service/schemas/events"
	"telemetry-service/suppliers"
)

type InitHouseService struct {
	houseRepository *repository.HouseRepository
	kafkaSupplier   suppliers.BrokerInterface
}

func NewInitHouseService(
	houseRepository *repository.HouseRepository,
	kafkaSupplier suppliers.BrokerInterface,
) *InitHouseService {
	return &InitHouseService{
		houseRepository: houseRepository,
		kafkaSupplier:   kafkaSupplier,
	}
}

func (s *InitHouseService) GetNewConnectedHouseEvent(ctx context.Context) (events.Event, error) {
	msg, err := s.kafkaSupplier.ReadNewHouseConnectedTopic(ctx)
	if err != nil {
		return events.Event{}, fmt.Errorf("failed to read message: %w", err)
	}

	var event events.Event
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return events.Event{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return event, nil
}

func (s *InitHouseService) ProcessEvent(event events.Event) error {
	data, ok := event.Payload.(events.InstallModuleToHousePayload)
	if !ok {
		return fmt.Errorf("invalid payload for InstallModuleToHouse event")
	}

	// some logic with data from event ...

	err := s.houseRepository.InsertHouse(event)
	if err != nil {
		return err
	}

	log.Println("msg data: ", data)
	return nil
}

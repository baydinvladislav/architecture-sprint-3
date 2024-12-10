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

type EmergencyService struct {
	deviceServiceSupplier *suppliers.DeviceServiceSupplier
	kafkaSupplier         *suppliers.KafkaSupplier
	emergencyRepository   *repository.EmergencyRepository
}

func NewEmergencyService(
	deviceServiceSupplier *suppliers.DeviceServiceSupplier,
	emergencyRepository *repository.EmergencyRepository,
	kafkaSupplier *suppliers.KafkaSupplier,
) *EmergencyService {
	return &EmergencyService{
		deviceServiceSupplier: deviceServiceSupplier,
		emergencyRepository:   emergencyRepository,
		kafkaSupplier:         kafkaSupplier,
	}
}

func (s *EmergencyService) GetEmergencyStopEvent(ctx context.Context) (events.BaseEvent, error) {
	msg, err := s.kafkaSupplier.ReadEmergencyStopTopic(ctx)
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

func (s *EmergencyService) ProcessEvent(event events.BaseEvent) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	var data events.EmergencyPayload
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		return fmt.Errorf("invalid payload for EmergencyShutdown event: %v", err)
	}

	log.Println("msg data finish log: ", data)

	err = s.emergencyRepository.SaveEmergencyEvent(event)
	if err != nil {
		return err
	}

	return nil
}

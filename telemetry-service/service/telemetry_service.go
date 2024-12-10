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

type TelemetryService struct {
	telemetryRepository *repository.TelemetryRepository
	kafkaSupplier       *suppliers.KafkaSupplier
}

func NewTelemetryService(
	telemetryRepository *repository.TelemetryRepository,
	kafkaSupplier *suppliers.KafkaSupplier,
) *TelemetryService {
	return &TelemetryService{
		telemetryRepository: telemetryRepository,
		kafkaSupplier:       kafkaSupplier,
	}
}

func (r *TelemetryService) GetTelemetryEvent(ctx context.Context) (events.BaseEvent, error) {
	msg, err := r.kafkaSupplier.ReadTelemetryTopic(ctx)
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

func (r *TelemetryService) ProcessEvent(event events.BaseEvent) error {
	var data events.TelemetryPayload

	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		return fmt.Errorf("failed to unmarshal payload to TelemetryPayload: %v", err)
	}

	log.Println("msg data: ", data)

	err = r.SaveEvent(event)
	if err != nil {
		return err
	}
	return nil
}

func (r *TelemetryService) SaveEvent(event events.BaseEvent) error {
	log.Println("Saving event...")

	err := r.telemetryRepository.InsertEvent(event)
	if err != nil {
		return fmt.Errorf("failed to save event: %v", err)
	}

	log.Println("BaseEvent saved successfully")
	return nil
}

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
	telemetryRepository repository.TelemetryRepositoryInterface
	kafkaSupplier       suppliers.BrokerInterface
}

func NewTelemetryService(
	telemetryRepository repository.TelemetryRepositoryInterface,
	kafkaSupplier suppliers.BrokerInterface,
) *TelemetryService {
	return &TelemetryService{
		telemetryRepository: telemetryRepository,
		kafkaSupplier:       kafkaSupplier,
	}
}

func (r *TelemetryService) GetTelemetryEvent(ctx context.Context) (events.Event, error) {
	msg, err := r.kafkaSupplier.ReadTelemetryTopic(ctx)
	if err != nil {
		return events.Event{}, fmt.Errorf("failed to read message: %w", err)
	}

	fmt.Printf("Raw Kafka message value: %s\n", string(msg.Value))

	var event events.Event
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return events.Event{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return event, nil
}

func (r *TelemetryService) ProcessEvent(event events.Event) error {
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

func (r *TelemetryService) SaveEvent(event events.Event) error {
	log.Println("Saving event...")

	err := r.telemetryRepository.InsertEvent(event)
	if err != nil {
		return fmt.Errorf("failed to save event: %v", err)
	}

	log.Println("Event saved successfully")
	return nil
}

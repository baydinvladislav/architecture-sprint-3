package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"telemetry-service/repository"
	"telemetry-service/schemas"
	"telemetry-service/suppliers"
)

type TelemetryService struct {
	kafkaSupplier         *suppliers.KafkaSupplier
	telemetryRepository   *repository.TelemetryRepository
	deviceServiceSupplier *suppliers.DeviceServiceSupplier
}

func NewTelemetryService(
	kafkaSupplier *suppliers.KafkaSupplier,
	deviceServiceSupplier *suppliers.DeviceServiceSupplier,
	telemetryRepository *repository.TelemetryRepository,
) *TelemetryService {
	return &TelemetryService{
		kafkaSupplier:         kafkaSupplier,
		telemetryRepository:   telemetryRepository,
		deviceServiceSupplier: deviceServiceSupplier,
	}
}

func (s *TelemetryService) ReadMessage(ctx context.Context) (schemas.Event, error) {
	msg, err := s.kafkaSupplier.ReadMessage(ctx)
	if err != nil {
		return schemas.Event{}, fmt.Errorf("failed to read message from Kafka: %v", err)
	}

	log.Printf("Received message: %s", string(msg.Value))

	var event schemas.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return schemas.Event{}, fmt.Errorf("failed to unmarshal event: %v", err)
	}

	return event, nil
}

func (s *TelemetryService) ProcessEvent(event schemas.Event) error {
	log.Println("ProcessEvent method is called (stub).")

	switch event.EventType {
	case "TelemetryData":
		data, ok := event.Payload.(schemas.TelemetryPayload)
		if !ok {
			return fmt.Errorf("invalid payload for TelemetryData event")
		}

		err := s.saveEvent(data)
		if err != nil {
			return err
		}

	case "EmergencyShutdown":
		data, ok := event.Payload.(schemas.EmergencyPayload)
		if !ok {
			return fmt.Errorf("invalid payload for EmergencyShutdown event")
		}

		err := s.saveEvent(data)
		if err != nil {
			return err
		}

		err = s.emergencyShutdownModule(data)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}

	return nil
}

func (s *TelemetryService) saveEvent(event schemas.EventPayload) error {
	log.Println("saveEvent method is called (stub).")
	return nil
}

func (s *TelemetryService) emergencyShutdownModule(event schemas.EmergencyPayload) error {
	log.Println("emergencyShutdownModule method is called (stub).")
	return nil
}

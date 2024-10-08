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
	log.Println("ProcessEvent method is called (stub).", event)

	switch event.EventType {
	case "TelemetryData":
		var data schemas.TelemetryPayload

		payloadBytes, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		if err := json.Unmarshal(payloadBytes, &data); err != nil {
			return fmt.Errorf("failed to unmarshal payload to TelemetryPayload: %v", err)
		}

		log.Println("msg data: ", data)

		err = s.saveEvent(data)
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

	case "InstallModuleToHouse":
		data, ok := event.Payload.(schemas.InstallModuleToHousePayload)
		if !ok {
			return fmt.Errorf("invalid payload for InstallModuleToHouse event")
		}

		err := s.saveEvent(data)
		if err != nil {
			return err
		}

		err = s.moduleConnect(data)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}

	return nil
}

func (s *TelemetryService) saveEvent(event schemas.EventPayload) error {
	log.Println("Saving event...")

	err := s.telemetryRepository.InsertEvent(event)
	if err != nil {
		return fmt.Errorf("failed to save event: %v", err)
	}

	log.Println("Event saved successfully")
	return nil
}

func (s *TelemetryService) emergencyShutdownModule(event schemas.EmergencyPayload) error {
	log.Printf("shutdown module %s because of reason %s\n", event.EquipmentID, event.Reason)
	return nil
}

func (s *TelemetryService) moduleConnect(event schemas.InstallModuleToHousePayload) error {
	log.Printf("new module %s to house %d", event.ModuleID, event.HouseID)
	return nil
}

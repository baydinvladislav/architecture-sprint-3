package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"telemetry-service/schemas"
	"telemetry-service/service"
	"telemetry-service/suppliers"
)

type KafkaDispatcher struct {
	telemetryService *service.TelemetryService
	emergencyService *service.EmergencyService
	initHouseService *service.InitHouseService
	kafkaSupplier    *suppliers.KafkaSupplier
}

func NewKafkaDispatcher(
	telemetryService *service.TelemetryService,
	emergencyService *service.EmergencyService,
	initHouseService *service.InitHouseService,
	kafkaSupplier *suppliers.KafkaSupplier,
) *KafkaDispatcher {
	return &KafkaDispatcher{
		telemetryService: telemetryService,
		emergencyService: emergencyService,
		initHouseService: initHouseService,
		kafkaSupplier:    kafkaSupplier,
	}
}

func (s *KafkaDispatcher) ReadMessage(ctx context.Context, topic string) (schemas.Event, error) {
	msg, err := s.kafkaSupplier.ReadMessage(ctx, topic)
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

func (s *KafkaDispatcher) ProcessEvent(event schemas.Event) error {
	log.Println("ProcessEvent method is called (stub).", event)

	// TODO: make map instead of switch
	switch event.EventType {
	case "TelemetryData":
		err := s.telemetryService.ProcessEvent(event)
		if err != nil {
			return err
		}

	case "EmergencyShutdown":
		err := s.emergencyService.ProcessEvent(event)
		if err != nil {
			return err
		}

	case "InstallModuleToHouse":
		err := s.initHouseService.ProcessEvent(event)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}

	return nil
}

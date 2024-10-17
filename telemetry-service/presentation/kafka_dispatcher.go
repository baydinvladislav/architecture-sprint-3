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

func (s *KafkaDispatcher) ReadMessage(ctx context.Context, topic string) error {
	msg, err := s.kafkaSupplier.ReadMessage(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to read message from Kafka: %v", err)
	}

	log.Printf("Received message from topic %s: %v", topic, msg)

	var event schemas.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	err = s.RouteEvent(event)
	if err != nil {
		return fmt.Errorf("failed to route event: %v", err)
	}

	return nil
}

func (s *KafkaDispatcher) RouteEvent(event schemas.Event) error {
	log.Println("RouteEvent method is called (stub).", event)

	eventHandlers := map[string]func(event schemas.Event) error{
		"TelemetryData":        s.telemetryService.ProcessEvent,
		"EmergencyShutdown":    s.emergencyService.ProcessEvent,
		"InstallModuleToHouse": s.initHouseService.ProcessEvent,
	}

	if handler, found := eventHandlers[event.EventType]; found {
		return handler(event)
	}

	log.Printf("Unknown event type: %s", event.EventType)
	return fmt.Errorf("unknown event type: %s", event.EventType)
}

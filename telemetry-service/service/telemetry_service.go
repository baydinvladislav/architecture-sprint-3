package service

import (
	"fmt"
	"log"
	"telemetry-service/repository"
	"telemetry-service/suppliers"
)

type TelemetryService struct {
	kafkaClient           *suppliers.KafkaClient
	telemetryRepository   *repository.TelemetryRepository
	deviceServiceSupplier *suppliers.DeviceServiceSupplier
}

func NewTelemetryService(
	kafkaClient *suppliers.KafkaClient,
	deviceServiceSupplier *suppliers.DeviceServiceSupplier,
	telemetryRepository *repository.TelemetryRepository,
) *TelemetryService {
	return &TelemetryService{
		kafkaClient:           kafkaClient,
		telemetryRepository:   telemetryRepository,
		deviceServiceSupplier: deviceServiceSupplier,
	}
}

func (s *TelemetryService) ReadTelemetryTopic() map[string]interface{} {
	log.Println("ReadTelemetryTopic method is called (stub).")
	return nil
}

func (s *TelemetryService) ProcessEvent(event map[string]interface{}) error {
	log.Println("ProcessEvent method is called (stub).")

	switch event.EventType {
	case "TelemetryData":
		data, ok := event.Payload.(models.SensorData)
		if !ok {
			return fmt.Errorf("invalid payload for TelemetryData event")
		}
		return container.TelemetryService.SendSensorData(data)

	case "EmergencyShutdown":
		eventData, ok := event.Payload.(models.EmergencyEvent)
		if !ok {
			return fmt.Errorf("invalid payload for EmergencyShutdown event")
		}
		return container.TelemetryService.EmergencyShutdownModule(eventData)

	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (s *TelemetryService) EmergencyShutdownModule() error {
	log.Println("EmergencyShutdownModule method is called (stub).")
	return nil
}

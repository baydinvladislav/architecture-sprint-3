package service

import (
	"encoding/json"
	"fmt"
	"log"
	"telemetry-service/repository"
	"telemetry-service/schemas"
)

type TelemetryService struct {
	telemetryRepository *repository.TelemetryRepository
}

func NewTelemetryService(telemetryRepository *repository.TelemetryRepository) *TelemetryService {
	return &TelemetryService{telemetryRepository: telemetryRepository}
}

func (s *TelemetryService) ProcessEvent(event schemas.Event) error {
	var data schemas.TelemetryPayload

	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		return fmt.Errorf("failed to unmarshal payload to TelemetryPayload: %v", err)
	}

	log.Println("msg data: ", data)

	err = s.SaveEvent(event)
	if err != nil {
		return err
	}
	return nil
}

func (s *TelemetryService) SaveEvent(event schemas.Event) error {
	log.Println("Saving event...")

	err := s.telemetryRepository.InsertEvent(event)
	if err != nil {
		return fmt.Errorf("failed to save event: %v", err)
	}

	log.Println("Event saved successfully")
	return nil
}

package service

import (
	"encoding/json"
	"fmt"
	"log"
	"telemetry-service/repository"
	"telemetry-service/schemas"
	"telemetry-service/suppliers"
)

type EmergencyService struct {
	deviceServiceSupplier *suppliers.DeviceServiceSupplier
	emergencyRepository   *repository.EmergencyRepository
}

func NewEmergencyService(
	deviceServiceSupplier *suppliers.DeviceServiceSupplier,
	emergencyRepository *repository.EmergencyRepository,
) *EmergencyService {
	return &EmergencyService{
		deviceServiceSupplier: deviceServiceSupplier,
		emergencyRepository:   emergencyRepository,
	}
}

func (s *EmergencyService) ProcessEvent(event schemas.Event) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	var data schemas.EmergencyPayload
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

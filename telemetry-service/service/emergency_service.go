package service

import (
	"fmt"
	"log"
	"telemetry-service/schemas"
	"telemetry-service/suppliers"
)

type EmergencyService struct {
	deviceServiceSupplier *suppliers.DeviceServiceSupplier
}

func NewEmergencyService(deviceServiceSupplier *suppliers.DeviceServiceSupplier) *EmergencyService {
	return &EmergencyService{deviceServiceSupplier: deviceServiceSupplier}
}

func (s *EmergencyService) ProcessEvent(event schemas.Event) error {
	data, ok := event.Payload.(schemas.EmergencyPayload)
	if !ok {
		return fmt.Errorf("invalid payload for EmergencyShutdown event")
	}

	log.Println("msg data: ", data)
	return nil
}

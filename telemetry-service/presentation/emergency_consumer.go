package presentation

import (
	"context"
	"log"
	"telemetry-service/service"
)

func StartEmergencyConsumer(ctx context.Context, service *service.EmergencyService) {
	go func() {
		for {
			event, err := service.GetEmergencyStopEvent(ctx)
			if err != nil {
				log.Printf("Error while reading EmergencyStopTopic: %v", err)
				continue
			}

			err = service.ProcessEvent(event)
			if err != nil {
				log.Printf("Error processing EmergencyStop event: %v", err)
				continue
			}

			log.Printf("Processed EmergencyStop event: %v", event)
		}
	}()
}

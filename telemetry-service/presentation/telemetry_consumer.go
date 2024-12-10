package presentation

import (
	"context"
	"log"
	"telemetry-service/service"
)

func StartTelemetryConsumer(ctx context.Context, service *service.TelemetryService) {
	go func() {
		for {
			event, err := service.GetTelemetryEvent(ctx)
			if err != nil {
				log.Printf("Error while reading TelemetryTopic: %v", err)
				continue
			}

			err = service.ProcessEvent(event)
			if err != nil {
				log.Printf("Error processing Telemetry event: %v", err)
				continue
			}

			log.Printf("Processed Telemetry event: %v", event)
		}
	}()
}

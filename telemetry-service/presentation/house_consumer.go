package presentation

import (
	"context"
	"log"
	"telemetry-service/service"
)

func StartHouseConsumer(ctx context.Context, service *service.InitHouseService) {
	go func() {
		for {
			event, err := service.GetNewConnectedHouseEvent(ctx)
			if err != nil {
				log.Printf("Error while reading NewConnectedHouseTopic: %v", err)
				continue
			}

			err = service.ProcessEvent(event)
			if err != nil {
				log.Printf("Error processing NewConnectedHouse event: %v", err)
				continue
			}

			log.Printf("Processed NewConnectedHouse event: %v", event)
		}
	}()
}

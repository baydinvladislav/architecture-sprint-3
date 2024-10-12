package presentation

import (
	"context"
	"log"
	"telemetry-service/shared"
)

func HandleTelemetryTopic(ctx context.Context, container *shared.AppContainer) {
	var topic = container.AppSettings.TelemetryTopic

	log.Printf("Starting listener for topic: %s", topic)

	for {
		event, err := container.KafkaDispatcher.ReadMessage(ctx, topic)
		if err != nil {
			log.Printf("Error while reading message from topic %s: %v", topic, err)
			continue
		}

		if err := container.KafkaDispatcher.RouteEvent(event); err != nil {
			log.Printf("Error handling event from topic %s: %v", topic, err)
		}
	}
}

func HandleEmergencyTopic(ctx context.Context, container *shared.AppContainer) {
	var topic = container.AppSettings.EmergencyTopic

	log.Printf("Starting listener for topic: %s", topic)

	for {
		event, err := container.KafkaDispatcher.ReadMessage(ctx, topic)
		if err != nil {
			log.Printf("Error while reading message from topic %s: %v", topic, err)
			continue
		}

		if err := container.KafkaDispatcher.RouteEvent(event); err != nil {
			log.Printf("Error handling event from topic %s: %v", topic, err)
		}
	}
}

func HandleNewHouseTopic(ctx context.Context, container *shared.AppContainer) {
	var topic = container.AppSettings.NewHouseTopic

	log.Printf("Starting listener for topic: %s", topic)

	for {
		event, err := container.KafkaDispatcher.ReadMessage(ctx, topic)
		if err != nil {
			log.Printf("Error while reading message from topic %s: %v", topic, err)
			continue
		}

		if err := container.KafkaDispatcher.RouteEvent(event); err != nil {
			log.Printf("Error handling event from topic %s: %v", topic, err)
		}
	}
}

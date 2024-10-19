package kafka_presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"user-service/schemas"
	"user-service/service"
	"user-service/suppliers"
)

type KafkaDispatcher struct {
	houseService  *service.HouseService
	kafkaSupplier *suppliers.KafkaSupplier
}

func NewKafkaDispatcher(
	kafkaSupplier *suppliers.KafkaSupplier,
	houseService *service.HouseService,
) *KafkaDispatcher {
	return &KafkaDispatcher{
		houseService:  houseService,
		kafkaSupplier: kafkaSupplier,
	}
}

func (s *KafkaDispatcher) ReadMessage(ctx context.Context, topic string) error {
	log.Printf("Read topic %s ...", topic)

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
	eventHandlers := map[string]func(event schemas.Event) error{}

	if handler, found := eventHandlers[event.EventType]; found {
		return handler(event)
	}

	log.Printf("Unknown event type: %s", event.EventType)
	return fmt.Errorf("unknown event type: %s", event.EventType)
}

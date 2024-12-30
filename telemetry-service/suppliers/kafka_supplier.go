package suppliers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSupplier struct {
	emergencyStopConsumer     *kafka.Reader
	newHouseConnectedConsumer *kafka.Reader
	telemetryConsumer         *kafka.Reader
}

func NewKafkaSupplier(
	kafkaBrokers []string,
	groupID string,
	emergencyStopTopic string,
	newHouseConnectedTopic string,
	telemetryTopic string,
) (*KafkaSupplier, error) {
	emergencyStopConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   emergencyStopTopic,
		GroupID: groupID,
	})

	newHouseConnectedConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   newHouseConnectedTopic,
		GroupID: groupID,
	})

	telemetryConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   telemetryTopic,
		GroupID: groupID,
	})

	return &KafkaSupplier{
		emergencyStopConsumer:     emergencyStopConsumer,
		newHouseConnectedConsumer: newHouseConnectedConsumer,
		telemetryConsumer:         telemetryConsumer,
	}, nil
}

func (ks *KafkaSupplier) ReadEmergencyStopTopic(ctx context.Context) (kafka.Message, error) {
	log.Printf("Listen topic EmergencyStopTopic...")

	msg, err := ks.emergencyStopConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message in EmergencyStopTopic: %v", msg)

	return msg, nil
}

func (ks *KafkaSupplier) ReadNewHouseConnectedTopic(ctx context.Context) (kafka.Message, error) {
	log.Printf("Listen topic NewHouseConnectedTopic...")

	msg, err := ks.newHouseConnectedConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message in NewHouseConnectedTopic: %v", msg)

	return msg, nil
}

func (ks *KafkaSupplier) ReadTelemetryTopic(ctx context.Context) (kafka.Message, error) {
	log.Printf("Listen topic TelemetryTopic...")

	msg, err := ks.telemetryConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message in TelemetryTopic: %v", msg)

	return msg, nil
}

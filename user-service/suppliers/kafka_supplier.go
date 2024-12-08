package suppliers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"user-service/schemas/events"
)

type KafkaSupplier struct {
	verifyHouseConsumer *kafka.Reader
	verifyHouseProducer *kafka.Writer
}

func NewKafkaSupplier(
	kafkaBroker string,
	moduleVerificationTopic string,
	moduleAdditionTopic string,
	groupID string,
) (*KafkaSupplier, error) {
	verifyHouseConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   moduleAdditionTopic,
		GroupID: groupID,
	})

	verifyHouseProducer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    moduleVerificationTopic,
		Balancer: &kafka.LeastBytes{},
	}

	ks := &KafkaSupplier{
		verifyHouseConsumer: verifyHouseConsumer,
		verifyHouseProducer: verifyHouseProducer,
	}

	return ks, nil
}

func (ks *KafkaSupplier) ReadModuleAdditionTopic(ctx context.Context) (kafka.Message, error) {
	log.Printf("Listen topic ModuleAdditionTopic")

	msg, err := ks.verifyHouseConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message in ModuleAdditionTopic: %v", msg)

	return msg, nil
}

func (ks *KafkaSupplier) SendMessageToModuleVerificationTopic(
	ctx context.Context,
	key []byte,
	event events.ModuleVerificationEvent,
) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	log.Printf(
		"Sending message to topic '%s': key=%s, value=%s",
		ks.verifyHouseProducer.Topic,
		string(key),
		string(value),
	)

	if err := ks.verifyHouseProducer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf(
		"Message successfully sent to topic '%s': key=%s, value=%s",
		ks.verifyHouseProducer.Topic,
		string(key),
		string(value),
	)

	return nil
}

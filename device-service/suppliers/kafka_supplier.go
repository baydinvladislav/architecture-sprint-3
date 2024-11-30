package suppliers

import (
	"context"
	"device-service/schemas/events"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSupplier struct {
	moduleAdditionProducer       *kafka.Writer
	equipmentChangeStateProducer *kafka.Writer
	moduleVerificationConsumer   *kafka.Reader
}

func NewKafkaSupplier(
	kafkaBroker string,
	moduleAdditionTopic string,
	moduleVerificationTopic string,
	equipmentChangeStateTopic string,
	groupID string,
) (*KafkaSupplier, error) {
	moduleAdditionProducer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    moduleAdditionTopic,
		Balancer: &kafka.LeastBytes{},
	}

	equipmentChangeStateProducer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    equipmentChangeStateTopic,
		Balancer: &kafka.LeastBytes{},
	}

	moduleVerificationConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   moduleVerificationTopic,
		GroupID: groupID,
	})

	ks := &KafkaSupplier{
		moduleAdditionProducer:       moduleAdditionProducer,
		moduleVerificationConsumer:   moduleVerificationConsumer,
		equipmentChangeStateProducer: equipmentChangeStateProducer,
	}

	return ks, nil
}

func (ks *KafkaSupplier) SendMessageToAdditionTopic(
	ctx context.Context,
	key []byte,
	event events.HomeVerificationEvent,
) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	err = ks.moduleAdditionProducer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (ks *KafkaSupplier) SendMessageToEquipmentChangeStateTopic(
	ctx context.Context,
	key []byte,
	event events.ChangeEquipmentStateEvent,
) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	if err := ks.equipmentChangeStateProducer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}

func (ks *KafkaSupplier) ReadModuleVerificationTopic(ctx context.Context) (kafka.Message, error) {
	log.Printf("Listen ModuleVerificationTopic ...")

	msg, err := ks.moduleVerificationConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message in ModuleVerificationTopic: %v", msg)

	return msg, nil
}

func (ks *KafkaSupplier) Close() {
	if err := ks.moduleAdditionProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka moduleVerificationProducer: %v", err)
	}

	if err := ks.equipmentChangeStateProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka equipmentChangeStateProducer: %v", err)
	}

	if err := ks.moduleVerificationConsumer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka moduleVerificationConsumer: %v", err)
	}
}

package suppliers

import (
	"context"
	"device-service/schemas"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSupplierInterface interface {
	SendMessageToAdditionTopic(ctx context.Context, key []byte, event schemas.HomeVerificationEvent) error
	SendMessageToEquipmentChangeStateTopic(ctx context.Context, key []byte, event schemas.ChangeEquipmentStateEvent) error
	ReadModuleVerificationTopic(ctx context.Context) (kafka.Message, error)
	Close()
}

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

	kc := &KafkaSupplier{
		moduleAdditionProducer:       moduleAdditionProducer,
		moduleVerificationConsumer:   moduleVerificationConsumer,
		equipmentChangeStateProducer: equipmentChangeStateProducer,
	}

	return kc, nil
}

func (kc *KafkaSupplier) SendMessageToAdditionTopic(
	ctx context.Context,
	key []byte,
	event schemas.HomeVerificationEvent,
) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	err = kc.moduleAdditionProducer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (kc *KafkaSupplier) SendMessageToEquipmentChangeStateTopic(
	ctx context.Context,
	key []byte,
	event schemas.ChangeEquipmentStateEvent,
) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	if err := kc.equipmentChangeStateProducer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}

func (kc *KafkaSupplier) ReadModuleVerificationTopic(ctx context.Context) (kafka.Message, error) {
	log.Printf("Read topicModuleVerificationTopic ...")

	msg, err := kc.moduleVerificationConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message: %v", msg)

	return msg, nil
}

func (kc *KafkaSupplier) Close() {
	if err := kc.moduleAdditionProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka moduleVerificationProducer: %v", err)
	}

	if err := kc.equipmentChangeStateProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka equipmentChangeStateProducer: %v", err)
	}

	if err := kc.moduleVerificationConsumer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka moduleVerificationConsumer: %v", err)
	}
}

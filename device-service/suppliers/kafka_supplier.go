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
	kafkaBrokers []string,
	moduleAdditionTopic string,
	moduleVerificationTopic string,
	equipmentChangeStateTopic string,
	groupID string,
) (*KafkaSupplier, error) {
	moduleAdditionProducer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers...),
		Topic:    moduleAdditionTopic,
		Balancer: &kafka.LeastBytes{},
	}

	equipmentChangeStateProducer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers...),
		Topic:    equipmentChangeStateTopic,
		Balancer: &kafka.LeastBytes{},
	}

	moduleVerificationConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
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

	log.Printf(
		"Sending message to topic '%s': key=%s, value=%s",
		ks.moduleAdditionProducer.Topic,
		string(key),
		string(value),
	)

	err = ks.moduleAdditionProducer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	log.Printf(
		"Message successfully sent to topic '%s': key=%s, value=%s",
		ks.moduleAdditionProducer.Topic,
		string(key),
		string(value),
	)

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

	log.Printf(
		"Sending message to topic '%s': key=%s, value=%s",
		ks.equipmentChangeStateProducer.Topic,
		string(key),
		string(value),
	)

	if err := ks.equipmentChangeStateProducer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf(
		"Message successfully sent to topic '%s': key=%s, value=%s",
		ks.equipmentChangeStateProducer.Topic,
		string(key),
		string(value),
	)

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

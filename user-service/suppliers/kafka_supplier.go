package suppliers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSupplier struct {
	moduleAdditionProducer     *kafka.Writer
	moduleVerificationConsumer *kafka.Reader
}

func NewKafkaSupplier(
	kafkaBroker string,
	moduleAdditionTopic string,
	moduleVerificationTopic string,
	groupID string,
) *KafkaSupplier {
	moduleAdditionProducer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    moduleVerificationTopic,
		Balancer: &kafka.LeastBytes{},
	}

	moduleVerificationConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   moduleAdditionTopic,
		GroupID: groupID,
	})

	return &KafkaSupplier{
		moduleAdditionProducer:     moduleAdditionProducer,
		moduleVerificationConsumer: moduleVerificationConsumer,
	}
}

func (kc *KafkaSupplier) SendMessage(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	err := kc.moduleAdditionProducer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (kc *KafkaSupplier) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := kc.moduleVerificationConsumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}
	return msg, nil
}

func (kc *KafkaSupplier) Close() {
	if err := kc.moduleAdditionProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka moduleVerificationProducer: %v", err)
	}

	if err := kc.moduleVerificationConsumer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka moduleVerificationConsumer: %v", err)
	}
}

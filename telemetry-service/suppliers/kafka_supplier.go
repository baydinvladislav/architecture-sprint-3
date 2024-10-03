package suppliers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSupplier struct {
	consumer *kafka.Reader
}

func NewKafkaSupplier(broker string, topic string, groupID string) *KafkaSupplier {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaSupplier{consumer: consumer}
}

func (kc *KafkaSupplier) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := kc.consumer.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}
	return msg, nil
}

func (kc *KafkaSupplier) SendMessage(msg kafka.Message) error {
	log.Printf("Sent message: %s", string(msg.Value))
	return nil
}

func (kc *KafkaSupplier) Close() {
	if err := kc.consumer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka consumer: %v", err)
	}
}

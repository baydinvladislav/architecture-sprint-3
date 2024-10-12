package suppliers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSupplier struct {
	brokers []string
	groupID string
}

func NewKafkaSupplier(brokers []string, groupID string) *KafkaSupplier {
	return &KafkaSupplier{
		brokers: brokers,
		groupID: groupID,
	}
}

func (kc *KafkaSupplier) createReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: kc.brokers,
		Topic:   topic,
		GroupID: kc.groupID,
	})
}

func (kc *KafkaSupplier) ReadMessage(ctx context.Context, topic string) (kafka.Message, error) {
	reader := kc.createReader(topic)
	defer reader.Close()

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	log.Printf("Received message from topic %s: %s", topic, string(msg.Value))
	return msg, nil
}

func (kc *KafkaSupplier) SendMessage(msg kafka.Message) error {
	log.Printf("Sent message: %s", string(msg.Value))
	return nil
}

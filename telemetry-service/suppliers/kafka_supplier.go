package suppliers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaClient struct {
	consumer *kafka.Reader
}

func NewKafkaClient(broker string, topic string, groupID string) *KafkaClient {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaClient{consumer: consumer}
}

func (kc *KafkaClient) Start() {
	for {
		m, err := kc.consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			continue
		}
		if err := kc.processMessage(m); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}
}

func (kc *KafkaClient) processMessage(msg kafka.Message) error {
	log.Printf("Received message: %s", string(msg.Value))

	// Здесь вы можете вызвать метод сервиса телеметрии для обработки сообщения
	// Например, если ваше сообщение - это данные датчика
	// вы можете сначала разобрать их в подходящую структуру, а затем вызвать метод сервиса

	// Пример для JSON
	//var sensorData shared.SensorData // Замените на вашу модель
	//if err := json.Unmarshal(msg.Value, &sensorData); err != nil {
	//	return err
	//}
	//
	return nil
}

func (kc *KafkaClient) Close() {
	if err := kc.consumer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka consumer: %v", err)
	}
}

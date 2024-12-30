package suppliers

import (
	"context"
	"device-service/schemas/events"
	"fmt"
	"github.com/segmentio/kafka-go"
)

var ErrKafkaSupplier = fmt.Errorf("error during send message in Kafka")

type KafkaSupplierInterface interface {
	SendMessageToAdditionTopic(ctx context.Context, key []byte, event events.HomeVerificationEvent) error
	SendMessageToEquipmentChangeStateTopic(ctx context.Context, key []byte, event events.ChangeEquipmentStateEvent) error
	ReadModuleVerificationTopic(ctx context.Context) (kafka.Message, error)
	Close()
}

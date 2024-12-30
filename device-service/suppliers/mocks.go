package suppliers

import (
	"context"
	"device-service/schemas/events"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type MockKafkaSupplier struct {
	mock.Mock
}

func (k *MockKafkaSupplier) ReadModuleVerificationTopic(ctx context.Context) (kafka.Message, error) {
	args := k.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (k *MockKafkaSupplier) SendMessageToEquipmentChangeStateTopic(
	ctx context.Context,
	key []byte,
	event events.ChangeEquipmentStateEvent,
) error {
	args := k.Called(ctx, key, event)
	return args.Error(0)
}

func (k *MockKafkaSupplier) SendMessageToAdditionTopic(
	ctx context.Context,
	key []byte,
	event events.HomeVerificationEvent,
) error {
	args := k.Called(ctx, key, event)
	return args.Error(0)
}

func (k *MockKafkaSupplier) Close() {
	k.Called()
}

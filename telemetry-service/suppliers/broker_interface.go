package suppliers

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type BrokerInterface interface {
	ReadEmergencyStopTopic(ctx context.Context) (kafka.Message, error)
	ReadNewHouseConnectedTopic(ctx context.Context) (kafka.Message, error)
	ReadTelemetryTopic(ctx context.Context) (kafka.Message, error)
}

package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"telemetry-service/repository"
	"telemetry-service/suppliers"
	"testing"
)

func TestTelemetryService_GetTelemetryEvent_Ok(t *testing.T) {
	// init tested code with mocks
	telemetryRepository := new(repository.MockTelemetryRepository)
	kafkaSupplier := new(suppliers.MockKafkaSupplier)
	telemetryService := NewTelemetryService(telemetryRepository, kafkaSupplier)

	// mock Kafka message with TelemetryPayload
	mockMessage := kafka.Message{
		Value: []byte(`{
			"event_type": "TelemetryEvent",
			"payload": {
				"source_id": "sensor_123",
				"source_type": "sensor",
				"value": 42.0,
				"time": 1672531200
			}}`),
	}

	kafkaSupplier.On("ReadTelemetryTopic", context.Background()).Return(mockMessage, nil)

	// call tested code
	event, err := telemetryService.GetTelemetryEvent(context.Background())

	// check no error
	require.NoError(t, err)

	// check event type
	require.Equal(t, "TelemetryEvent", event.EventType)

	// check payload
	payload, ok := event.Payload.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "sensor_123", payload["source_id"])
	require.Equal(t, "sensor", payload["source_type"])
	require.Equal(t, 42.0, payload["value"])
	require.Equal(t, float64(1672531200), payload["time"])
}

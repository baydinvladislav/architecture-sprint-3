package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"telemetry-service/repository"
	"telemetry-service/schemas/events"
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
	payload, ok := event.Payload.(events.TelemetryPayload)
	require.True(t, ok)
	require.Equal(t, "sensor_123", payload.SourceID)
	require.Equal(t, events.SourceTypeSensor, payload.SourceType)
	require.Equal(t, 42.0, payload.Value)
	require.Equal(t, int64(1672531200), payload.Time)
}

//func TestTelemetryService_GetTelemetryEvent_ReadError(t *testing.T) {
//	// init tested code with mocks
//	telemetryRepository := new(repository.MockTelemetryRepository)
//	kafkaSupplier := new(suppliers.MockKafkaSupplier)
//	telemetryService := NewTelemetryService(telemetryRepository, kafkaSupplier)
//
//	// mock Kafka read error
//	kafkaSupplier.On("ReadTelemetryTopic", context.Background()).Return(nil, errors.New("read error"))
//
//	// call tested code
//	event, err := telemetryService.GetTelemetryEvent(context.Background())
//
//	// check error
//	require.Error(t, err)
//	require.Contains(t, err.Error(), "read error")
//
//	// check empty event returned
//	require.Equal(t, events.Event{}, event)
//}
//
//func TestTelemetryService_ProcessEvent_TelemetryPayload_Ok(t *testing.T) {
//	// init tested code with mocks
//	telemetryRepository := new(repository.MockTelemetryRepository)
//	kafkaSupplier := new(suppliers.MockKafkaSupplier)
//	telemetryService := NewTelemetryService(telemetryRepository, kafkaSupplier)
//
//	// prepare TelemetryPayload event
//	event := events.Event{
//		EventType: "TelemetryEvent",
//		Payload: events.TelemetryPayload{
//			SourceID:   "sensor_123",
//			SourceType: events.SourceTypeSensor,
//			Value:      42.0,
//			Time:       1672531200,
//		},
//	}
//
//	// mock repository save success
//	telemetryRepository.On("InsertEvent", event).Return(nil)
//
//	// call tested code
//	err := telemetryService.ProcessEvent(event)
//
//	// check no error
//	require.NoError(t, err)
//
//	// check repository save called
//	telemetryRepository.AssertCalled(t, "InsertEvent", event)
//}
//
//func TestTelemetryService_ProcessEvent_EmergencyPayload_Ok(t *testing.T) {
//	// init tested code with mocks
//	telemetryRepository := new(repository.MockTelemetryRepository)
//	kafkaSupplier := new(suppliers.MockKafkaSupplier)
//	telemetryService := NewTelemetryService(telemetryRepository, kafkaSupplier)
//
//	// prepare EmergencyPayload event
//	event := events.Event{
//		EventType: "EmergencyEvent",
//		Payload: events.EmergencyPayload{
//			SourceID: "sensor_123",
//			Reason:   "Overheat",
//		},
//	}
//
//	// mock repository save success
//	telemetryRepository.On("InsertEvent", event).Return(nil)
//
//	// call tested code
//	err := telemetryService.ProcessEvent(event)
//
//	// check no error
//	require.NoError(t, err)
//
//	// check repository save called
//	telemetryRepository.AssertCalled(t, "InsertEvent", event)
//}
//
//func TestTelemetryService_SaveEvent_TelemetryPayload_Ok(t *testing.T) {
//	// init tested code with mocks
//	telemetryRepository := new(repository.MockTelemetryRepository)
//	kafkaSupplier := new(suppliers.MockKafkaSupplier)
//	telemetryService := NewTelemetryService(telemetryRepository, kafkaSupplier)
//
//	// prepare TelemetryPayload event
//	event := events.Event{
//		EventType: "TelemetryEvent",
//		Payload: events.TelemetryPayload{
//			SourceID:   "sensor_123",
//			SourceType: events.SourceTypeSensor,
//			Value:      42.0,
//			Time:       1672531200,
//		},
//	}
//
//	// mock repository save success
//	telemetryRepository.On("InsertEvent", event).Return(nil)
//
//	// call tested code
//	err := telemetryService.SaveEvent(event)
//
//	// check no error
//	require.NoError(t, err)
//
//	// check repository save called
//	telemetryRepository.AssertCalled(t, "InsertEvent", event)
//}
//
//func TestTelemetryService_SaveEvent_Error(t *testing.T) {
//	// init tested code with mocks
//	telemetryRepository := new(repository.MockTelemetryRepository)
//	kafkaSupplier := new(suppliers.MockKafkaSupplier)
//	telemetryService := NewTelemetryService(telemetryRepository, kafkaSupplier)
//
//	// prepare EmergencyPayload event
//	event := events.Event{
//		EventType: "EmergencyEvent",
//		Payload: events.EmergencyPayload{
//			SourceID: "sensor_123",
//			Reason:   "Overheat",
//		},
//	}
//
//	// mock repository save error
//	telemetryRepository.On("InsertEvent", event).Return(errors.New("save error"))
//
//	// call tested code
//	err := telemetryService.SaveEvent(event)
//
//	// check error
//	require.Error(t, err)
//	require.Contains(t, err.Error(), "save error")
//}

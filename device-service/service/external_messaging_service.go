package service

import (
	"context"
	"device-service/schemas/events"
	"device-service/suppliers"
	"encoding/json"
	"fmt"
)

type ExternalMessagingService struct {
	supplier suppliers.KafkaSupplierInterface
}

func NewExternalMessagingService(supplier suppliers.KafkaSupplierInterface) *ExternalMessagingService {
	return &ExternalMessagingService{
		supplier: supplier,
	}
}

func (s *ExternalMessagingService) SendModuleAdditionEvent(
	ctx context.Context,
	key []byte,
	event events.HomeVerificationEvent,
) error {
	return s.supplier.SendMessageToAdditionTopic(ctx, key, event)
}

func (s *ExternalMessagingService) SendEquipmentStateChangeEvent(
	ctx context.Context,
	key []byte,
	event events.ChangeEquipmentStateEvent,
) error {
	return s.supplier.SendMessageToEquipmentChangeStateTopic(ctx, key, event)
}

func (s *ExternalMessagingService) ReadModuleVerificationEvent(
	ctx context.Context,
) (events.BaseEvent, error) {
	msg, err := s.supplier.ReadModuleVerificationTopic(ctx)
	if err != nil {
		return events.BaseEvent{}, fmt.Errorf("failed to read message: %w", err)
	}

	var event events.BaseEvent
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return events.BaseEvent{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return event, nil
}

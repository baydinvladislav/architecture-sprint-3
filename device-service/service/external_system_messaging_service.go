package service

import (
	"context"
	"device-service/schemas"
	"device-service/suppliers"
	"encoding/json"
	"fmt"
)

type ExternalSystemMessagingService struct {
	supplier suppliers.KafkaSupplierInterface
}

func NewExternalSystemMessagingService(supplier suppliers.KafkaSupplierInterface) *ExternalSystemMessagingService {
	return &ExternalSystemMessagingService{
		supplier: supplier,
	}
}

func (s *ExternalSystemMessagingService) SendModuleAdditionEvent(
	ctx context.Context,
	key []byte,
	event schemas.HomeVerificationEvent,
) error {
	return s.supplier.SendMessageToAdditionTopic(ctx, key, event)
}

func (s *ExternalSystemMessagingService) SendEquipmentStateChangeEvent(
	ctx context.Context,
	key []byte,
	event schemas.ChangeEquipmentStateEvent,
) error {
	return s.supplier.SendMessageToEquipmentChangeStateTopic(ctx, key, event)
}

func (s *ExternalSystemMessagingService) ReadModuleVerificationEvent(
	ctx context.Context,
) (schemas.BaseEvent, error) {
	msg, err := s.supplier.ReadModuleVerificationTopic(ctx)
	if err != nil {
		return schemas.BaseEvent{}, fmt.Errorf("failed to read message: %w", err)
	}

	var event schemas.BaseEvent
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		return schemas.BaseEvent{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return event, nil
}

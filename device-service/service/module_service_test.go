package service

import (
	"device-service/repository"
	"device-service/schemas/events"
	"device-service/suppliers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessMessage_Accepted(t *testing.T) {
	moduleRepository := new(repository.MockModuleRepository)
	kafkaSupplier := new(suppliers.MockKafkaSupplier)

	persistenceService := NewModulePersistenceService(moduleRepository)
	messagingService := NewExternalMessagingService(kafkaSupplier)

	moduleService := NewModuleService(persistenceService, messagingService)

	houseID := uuid.New()
	moduleID := uuid.New()

	event := events.BaseEvent{
		EventType: "ModuleVerificationEvent",
		Payload: events.ModuleVerificationEvent{
			HouseID:  houseID.String(),
			ModuleID: moduleID.String(),
			Decision: "ACCEPTED",
		},
	}

	moduleRepository.On("AcceptAdditionModuleToHouse", houseID, moduleID).Return(nil)

	err := moduleService.ProcessModuleVerificationEvent(event)

	require.NoError(t, err)
	moduleRepository.AssertCalled(t, "AcceptAdditionModuleToHouse", houseID, moduleID)
}

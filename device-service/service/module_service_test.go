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
	// init tested code with mocks
	moduleRepository := new(repository.MockModuleRepository)
	kafkaSupplier := new(suppliers.MockKafkaSupplier)
	persistenceService := NewModulePersistenceService(moduleRepository)
	messagingService := NewExternalMessagingService(kafkaSupplier)
	moduleService := NewModuleService(persistenceService, messagingService)

	// we got event from UserService with ACCEPTED decision
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

	// mock success inserting in database
	moduleRepository.On("AcceptAdditionModuleToHouse", houseID, moduleID).Return(nil)

	// call tested code
	err := moduleService.ProcessModuleVerificationEvent(event)

	// check calling db repository method
	moduleRepository.AssertCalled(t, "AcceptAdditionModuleToHouse", houseID, moduleID)

	// check no error
	require.NoError(t, err)
}

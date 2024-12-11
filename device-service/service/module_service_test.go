package service

import (
	"device-service/repository"
	"device-service/schemas/dto"
	"device-service/schemas/events"
	"device-service/suppliers"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleService_GetAllModules_Ok(t *testing.T) {
	// init tested code with mocks
	moduleRepository := new(repository.MockModuleRepository)
	kafkaSupplier := new(suppliers.MockKafkaSupplier)
	persistenceService := NewModulePersistenceService(moduleRepository)
	messagingService := NewExternalMessagingService(kafkaSupplier)
	moduleService := NewModuleService(persistenceService, messagingService)

	// mocks some records in database which repository returns
	testState1, _ := json.Marshal(map[string]interface{}{"running": "on"})
	testState2, _ := json.Marshal(map[string]interface{}{"running": "off"})
	mockModules := []dto.ModuleDto{
		{ID: uuid.New(), Type: "Type1", Description: "Module 1", State: string(testState1)},
		{ID: uuid.New(), Type: "Type2", Description: "Module 2", State: string(testState2)},
	}
	moduleRepository.On("GetAllModules").Return(mockModules, nil)

	// call tested code
	modules, err := moduleService.GetAllModules()

	// check no error
	require.NoError(t, err)

	// check that we got the same modules that we mocked from database
	require.Len(t, modules, 2)
	require.Equal(t, mockModules[0].ID, modules[0].ID)
	require.Equal(t, mockModules[1].Description, modules[1].Description)
}

func TestModuleService_ProcessModuleVerificationEvent_Accepted(t *testing.T) {
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

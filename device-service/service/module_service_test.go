package service

import (
	"context"
	"device-service/schemas/events"
	web_schemas "device-service/schemas/web"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockModuleRepository struct {
	mock.Mock
}

func (m *MockModuleRepository) GetAllModules() ([]web_schemas.ModuleOut, error) {
	args := m.Called()
	return args.Get(0).([]web_schemas.ModuleOut), args.Error(1)
}

func (m *MockModuleRepository) GetModulesByHouseID(houseID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	args := m.Called(houseID)
	return args.Get(0).([]web_schemas.ModuleOut), args.Error(1)
}

func (m *MockModuleRepository) TurnOnModule(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) TurnOffModule(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) GetModuleState(houseID, moduleID uuid.UUID) (*web_schemas.HouseModuleState, error) {
	args := m.Called(houseID, moduleID)
	return args.Get(0).(*web_schemas.HouseModuleState), args.Error(1)
}

func (m *MockModuleRepository) AcceptAdditionModuleToHouse(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) FailAdditionModuleToHouse(houseID, moduleID uuid.UUID) error {
	args := m.Called(houseID, moduleID)
	return args.Error(0)
}

func (m *MockModuleRepository) RequestAddingModuleToHouse(houseID uuid.UUID, moduleID uuid.UUID) ([]web_schemas.ModuleOut, error) {
	args := m.Called(houseID, moduleID)
	return args.Get(0).([]web_schemas.ModuleOut), args.Error(1)
}

func (m *MockModuleRepository) InsertNewHouseModuleState(houseModuleId uuid.UUID, state map[string]interface{}) error {
	args := m.Called(houseModuleId, state)
	return args.Error(0)
}

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

func TestProcessMessage_Accepted(t *testing.T) {
	moduleRepository := new(MockModuleRepository)
	kafkaSupplier := new(MockKafkaSupplier)

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

	success, err := moduleService.ProcessMessage(event)

	require.NoError(t, err)
	require.True(t, success)
	moduleRepository.AssertCalled(t, "AcceptAdditionModuleToHouse", houseID, moduleID)
}

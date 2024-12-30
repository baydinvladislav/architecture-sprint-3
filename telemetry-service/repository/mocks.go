package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockTelemetryRepository struct {
	mock.Mock
}

func (m *MockTelemetryRepository) InsertEvent(event interface{}) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockTelemetryRepository) Close() {
	m.Called()
}

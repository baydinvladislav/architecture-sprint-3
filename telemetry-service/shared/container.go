package shared

import (
	"context"
	"telemetry-service/presentation"
	"telemetry-service/repository"
	"telemetry-service/service"
	"telemetry-service/suppliers"
)

type AppContainer struct {
	KafkaDispatcher  *presentation.KafkaDispatcher
	AppSettings      *AppSettings
	TelemetryService *service.TelemetryService
	EmergencyService *service.EmergencyService
	InitHouseService *service.InitHouseService
}

func NewAppContainer(ctx context.Context) *AppContainer {
	appSettings := NewAppSettings()

	telemetryRepository := repository.NewTelemetryRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)

	houseRepository := repository.NewHouseRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)

	kafkaSupplier := suppliers.NewKafkaSupplier(
		[]string{appSettings.KafkaBroker},
		appSettings.GroupID,
	)

	deviceServiceSupplier := suppliers.NewDeviceServiceSupplier(appSettings.DeviceServiceUrl)

	telemetryService := service.NewTelemetryService(telemetryRepository)
	emergencyService := service.NewEmergencyService(deviceServiceSupplier)
	initHouseService := service.NewInitHouseService(houseRepository)

	kafkaDispatcher := presentation.NewKafkaDispatcher(
		telemetryService,
		emergencyService,
		initHouseService,
		kafkaSupplier,
	)
	return &AppContainer{
		TelemetryService: telemetryService,
		EmergencyService: emergencyService,
		InitHouseService: initHouseService,
		KafkaDispatcher:  kafkaDispatcher,
		AppSettings:      appSettings,
	}
}

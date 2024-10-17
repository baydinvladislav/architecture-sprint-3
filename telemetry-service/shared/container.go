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
	telemetryService := service.NewTelemetryService(telemetryRepository)

	deviceServiceSupplier := suppliers.NewDeviceServiceSupplier(appSettings.DeviceServiceUrl)

	emergencyRepository := repository.NewEmergencyRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)
	emergencyService := service.NewEmergencyService(deviceServiceSupplier, emergencyRepository)

	houseRepository := repository.NewHouseRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)
	initHouseService := service.NewInitHouseService(houseRepository)

	kafkaSupplier := suppliers.NewKafkaSupplier(
		[]string{appSettings.KafkaBroker},
		appSettings.GroupID,
	)

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

package shared

import (
	"context"
	"log"
	"telemetry-service/repository"
	"telemetry-service/service"
	"telemetry-service/suppliers"
)

type AppContainer struct {
	AppSettings      *AppSettings
	TelemetryService *service.TelemetryService
	EmergencyService *service.EmergencyService
	InitHouseService *service.InitHouseService
}

func NewAppContainer(ctx context.Context) *AppContainer {
	appSettings := NewAppSettings()

	kafkaSupplier, err := suppliers.NewKafkaSupplier(
		appSettings.KafkaBrokers,
		appSettings.GroupID,
		appSettings.EmergencyStopTopic,
		appSettings.NewHouseConnectedTopic,
		appSettings.TelemetryTopic,
	)
	if err != nil {
		log.Fatalf("failed to initialize Kafka supplier: %v", err)
	}

	deviceServiceSupplier := suppliers.NewDeviceServiceSupplier(appSettings.DeviceServiceUrl)

	telemetryRepository := repository.NewTelemetryRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)
	telemetryService := service.NewTelemetryService(
		telemetryRepository,
		kafkaSupplier,
	)

	emergencyRepository := repository.NewEmergencyRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)
	emergencyService := service.NewEmergencyService(
		deviceServiceSupplier,
		emergencyRepository,
		kafkaSupplier,
	)

	houseRepository := repository.NewHouseRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.TelemetryCollection,
	)
	initHouseService := service.NewInitHouseService(
		houseRepository,
		kafkaSupplier,
	)

	return &AppContainer{
		TelemetryService: telemetryService,
		EmergencyService: emergencyService,
		InitHouseService: initHouseService,
		AppSettings:      appSettings,
	}
}

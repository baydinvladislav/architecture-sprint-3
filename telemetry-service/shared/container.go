package shared

import (
	"context"
	"telemetry-service/repository"
	"telemetry-service/service"
	"telemetry-service/suppliers"
)

type AppContainer struct {
	TelemetryService *service.TelemetryService
	AppSettings      *AppSettings
}

func NewAppContainer(ctx context.Context) *AppContainer {
	appSettings := NewAppSettings()

	telemetryRepo := repository.NewTelemetryRepository(
		appSettings.MongoURI,
		appSettings.DatabaseName,
		appSettings.CollectionName,
	)

	kafkaClient := suppliers.NewKafkaClient(
		appSettings.KafkaBroker,
		appSettings.TelemetryTopic,
		appSettings.GroupID,
	)

	deviceServiceSupplier := suppliers.NewDeviceServiceSupplier(appSettings.DeviceServiceUrl)

	telemetryService := service.NewTelemetryService(
		kafkaClient,
		deviceServiceSupplier,
		telemetryRepo,
	)
	return &AppContainer{
		TelemetryService: telemetryService,
		AppSettings:      appSettings,
	}
}

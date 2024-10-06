package shared

import (
	"context"
	"device-service/persistance"
	"device-service/repository"
	"device-service/service"
	"device-service/suppliers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Container struct {
	ModuleService *service.ModuleService
	AppSettings   *AppSettings
	KafkaSupplier *suppliers.KafkaSupplier
}

func NewAppContainer(ctx context.Context) *Container {
	appSettings := NewAppSettings()
	dsn := appSettings.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	persistance.Migrate(db)

	moduleRepo := repository.NewGORMModuleRepository(db)
	moduleService := service.NewModuleService(moduleRepo)

	kafkaSupplier := suppliers.NewKafkaSupplier(
		appSettings.KafkaBroker,
		appSettings.ModuleAddedKafkaTopic,
		appSettings.ModuleVerificationKafkaTopic,
		appSettings.KafkaGroupID,
	)

	return &Container{
		ModuleService: moduleService,
		AppSettings:   appSettings,
		KafkaSupplier: kafkaSupplier,
	}
}

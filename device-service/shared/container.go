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
}

func NewAppContainer(ctx context.Context) *Container {
	log.Printf("Starting initializing application container")

	appSettings := NewAppSettings()
	dsn := appSettings.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	persistance.Migrate(db)

	moduleRepository := repository.NewGORMModuleRepository(db)
	persistenceService := service.NewModulePersistenceService(moduleRepository)

	kafkaSupplier, err := suppliers.NewKafkaSupplier(
		appSettings.KafkaBroker,
		appSettings.ModuleAddedKafkaTopic,
		appSettings.ModuleVerificationKafkaTopic,
		appSettings.EquipmentChangeStateTopic,
		appSettings.KafkaGroupID,
	)
	messagingService := service.NewExternalMessagingService(kafkaSupplier)
	if err != nil {
		log.Fatalf("Error initializing KafkaSupplier: %v", err)
	}

	moduleService := service.NewModuleService(
		persistenceService,
		messagingService,
	)

	appContainer := &Container{
		ModuleService: moduleService,
		AppSettings:   appSettings,
	}

	log.Printf("Application container successfully initialized")

	return appContainer
}

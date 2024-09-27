package shared

import (
	"context"
	"device-service/persistance"
	"device-service/repository"
	"device-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Container struct {
	ModuleService *service.ModuleService
	AppSettings   *AppSettings
}

func NewAppContainer(ctx context.Context) *Container {
	settings := NewAppSettings()
	dsn := settings.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&persistance.ModuleModel{}, persistance.HouseModuleModel{}, persistance.Device{})
	if err != nil {
		return nil
	}

	moduleRepo := repository.NewGORMModuleRepository(db)
	moduleService := service.NewModuleService(moduleRepo)
	return &Container{
		ModuleService: moduleService,
	}
}

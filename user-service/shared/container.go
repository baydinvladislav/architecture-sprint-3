package shared

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-service/persistance"
	"user-service/repository"
	"user-service/service"
	"user-service/suppliers"
)

type Container struct {
	UserService  *service.UserService
	HouseService *service.HouseService
	AppSettings  *AppSettings
}

func NewAppContainer(ctx context.Context) *Container {
	appSettings := NewAppSettings()
	dsn := appSettings.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&persistance.UserModel{}, &persistance.HouseModel{})
	if err != nil {
		log.Fatalf("failed to auto-migrate database models: %v", err)
		return nil
	}

	var accessSecret = []byte("+AAlQmR/sSml0D0QgZ9suJZwtLxHbJAzjvWLYsiER+0=")
	var refreshSecret = []byte("2hXKd7hDB/28TBKPyR262qVfDi1aX2t00IG99q6wxEc=")
	authService := service.NewAuthService(accessSecret, refreshSecret)

	userRepository := repository.NewGORMUserRepository(db)
	userService := service.NewUserService(authService, userRepository)

	kafkaSupplier, err := suppliers.NewKafkaSupplier(
		appSettings.KafkaBroker,
		appSettings.ModuleVerificationTopic,
		appSettings.ModuleAdditionTopic,
		appSettings.KafkaGroupID,
	)
	if err != nil {
		log.Fatalf("failed to initialize Kafka supplier: %v", err)
	}

	verifyService := service.NewVerifyConnectionService(
		appSettings.MinHomeSquare,
		appSettings.MaxHomeSquare,
	)

	houseRepository := repository.NewGORMHouseRepository(db)
	houseService := service.NewHouseService(
		houseRepository,
		userService,
		verifyService,
		kafkaSupplier,
	)

	return &Container{
		UserService:  userService,
		HouseService: houseService,
		AppSettings:  appSettings,
	}
}

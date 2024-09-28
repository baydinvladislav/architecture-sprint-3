package shared

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-service/persistance"
	"user-service/repository"
	"user-service/service"
)

type Container struct {
	AuthService  *service.AuthService
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
		return nil
	}

	userRepo := repository.NewGORMUserRepository(db)
	userService := service.NewUserService(userRepo)

	var accessSecret = []byte("+AAlQmR/sSml0D0QgZ9suJZwtLxHbJAzjvWLYsiER+0=")
	var refreshSecret = []byte("2hXKd7hDB/28TBKPyR262qVfDi1aX2t00IG99q6wxEc=")
	authService := service.NewAuthService(accessSecret, refreshSecret)

	houseRepo := repository.NewGORMHouseRepository(db)
	houseService := service.NewHouseService(houseRepo)
	return &Container{
		UserService:  userService,
		AuthService:  authService,
		HouseService: houseService,
		AppSettings:  appSettings,
	}
}

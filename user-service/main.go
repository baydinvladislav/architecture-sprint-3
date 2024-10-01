package main

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "user-service/docs"
	"user-service/middleware"
	"user-service/presentation"
	"user-service/shared"
)

func CreateApp(ctx context.Context) *gin.Engine {
	appContainer := shared.NewAppContainer(ctx)
	r := gin.Default()

	serviceGroup := r.Group("/service")
	{
		serviceGroup.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "OK"}) })
		serviceGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.POST("/register", func(c *gin.Context) { presentation.RegisterUser(c, appContainer) })
	r.POST("/login", func(c *gin.Context) { presentation.LoginUser(c, appContainer) })
	r.POST("/refresh-token", func(c *gin.Context) { presentation.RefreshToken(c, appContainer) })

	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware(appContainer.AuthService))
	{
		authGroup.GET("/users/:userId", func(c *gin.Context) { presentation.GetUserById(c, appContainer) })
		authGroup.POST("/houses", func(c *gin.Context) { presentation.CreateUserHouse(c, appContainer) })
		authGroup.GET("/houses", func(c *gin.Context) { presentation.GetUserHouses(c, appContainer) })
		authGroup.PUT("/houses/:houseId", func(c *gin.Context) { presentation.UpdateUserHouse(c, appContainer) })
	}

	return r
}

func main() {
	ctx := context.Background()
	app := CreateApp(ctx)

	err := app.Run(":8080")
	if err != nil {
		return
	}
}

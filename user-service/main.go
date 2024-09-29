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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", func(c *gin.Context) { presentation.RegisterUser(c, appContainer) })
		userGroup.POST("/login", func(c *gin.Context) { presentation.LoginUser(c, appContainer) })
		userGroup.POST("/refresh-token", func(c *gin.Context) { presentation.RefreshToken(c, appContainer) })

		userGroup.Use(middleware.AuthMiddleware(appContainer.AuthService))
		userGroup.GET("/:userId", func(c *gin.Context) { presentation.GetUserById(c, appContainer) })
	}

	houseGroup := r.Group("/houses")
	{
		userGroup.Use(middleware.AuthMiddleware(appContainer.AuthService))
		houseGroup.POST("/", func(c *gin.Context) { presentation.CreateUserHouse(c, appContainer) })
		houseGroup.GET("/", func(c *gin.Context) { presentation.GetUserHouses(c, appContainer) })
		houseGroup.PUT("/:houseId", func(c *gin.Context) { presentation.UpdateUserHouse(c, appContainer) })
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

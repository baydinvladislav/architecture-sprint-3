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

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", func(c *gin.Context) { presentation.RegisterUser(c, appContainer) })
	r.POST("/login", func(c *gin.Context) { presentation.LoginUser(c, appContainer) })
	r.POST("/refresh-token", func(c *gin.Context) { presentation.RefreshToken(c, appContainer) })

	r.Use(middleware.AuthMiddleware(appContainer.AuthService))
	r.GET("/users/:userId", func(c *gin.Context) { presentation.GetUserById(c, appContainer) })

	r.POST("/houses", func(c *gin.Context) { presentation.CreateUserHouse(c, appContainer) })
	r.GET("/houses", func(c *gin.Context) { presentation.GetUserHouses(c, appContainer) })
	r.PUT("/houses/:houseId", func(c *gin.Context) { presentation.UpdateUserHouse(c, appContainer) })

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

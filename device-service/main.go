package main

import (
	"context"
	"device-service/presentation"
	"device-service/shared"
	"github.com/gin-gonic/gin"
)

func CreateApp(ctx context.Context) *gin.Engine {
	appContainer := shared.NewAppContainer(ctx)
	r := gin.Default()

	moduleGroup := r.Group("/modules")
	{
		moduleGroup.GET("/", func(c *gin.Context) { presentation.GetAvailableModules(c, appContainer) })
	}

	return r
}

func main() {
	ctx := context.Background()
	app := CreateApp(ctx)

	err := app.Run(":8081")
	if err != nil {
		return
	}
}

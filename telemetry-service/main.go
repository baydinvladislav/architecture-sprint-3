package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"telemetry-service/shared"
)

func CreateApp(ctx context.Context, container *shared.AppContainer) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	go func() {
		for {
			event := container.TelemetryService.ReadTelemetryTopic()
			if err := container.TelemetryService.ProcessEvent(event); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		}
	}()

	return r
}

func main() {
	ctx := context.Background()
	appContainer := shared.NewAppContainer(ctx)

	app := CreateApp(ctx, appContainer)

	err := app.Run(":8082")
	if err != nil {
		log.Fatalf("Failed to run the app: %v", err)
		return
	}
}

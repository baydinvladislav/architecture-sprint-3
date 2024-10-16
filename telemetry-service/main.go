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

	initKafkaHandlers(ctx, container)

	return r
}

func initKafkaHandlers(ctx context.Context, container *shared.AppContainer) {
	go handleKafkaTopic(ctx, container, "telemetry.data")
	go handleKafkaTopic(ctx, container, "forced.module.shutdown")
}

func handleKafkaTopic(ctx context.Context, container *shared.AppContainer, topic string) {
	log.Printf("Starting Kafka consumer for topic: %s", topic)

	for {
		err := container.KafkaDispatcher.ReadMessage(ctx, topic)
		if err != nil {
			log.Printf("Error while reading message from topic %s: %v", topic, err)
			continue
		}
	}
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

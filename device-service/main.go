package main

import (
	"context"
	"device-service/presentation"
	"device-service/shared"
	"github.com/gin-gonic/gin"
	"log"
)

func CreateApp(ctx context.Context) *gin.Engine {
	appContainer := shared.NewAppContainer(ctx)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	moduleGroup := r.Group("/modules")
	{
		moduleGroup.GET(
			"/",
			func(c *gin.Context) { presentation.GetAvailableModules(c, appContainer) },
		)
		moduleGroup.GET(
			"/houses/:houseID",
			func(c *gin.Context) { presentation.GetModulesByHouseId(c, appContainer) },
		)
		moduleGroup.POST(
			"/houses/:houseID/modules/:moduleID/assign",
			func(c *gin.Context) { presentation.AddModuleToHouse(c, appContainer) },
		)
		moduleGroup.GET(
			"/houses/:houseID/modules/:moduleID/state",
			func(c *gin.Context) { presentation.GetModuleState(c, appContainer) },
		)
		moduleGroup.POST(
			"/houses/:houseID/modules/:moduleID/state",
			func(c *gin.Context) { presentation.ChangeModuleState(c, appContainer) },
		)
		moduleGroup.PUT(
			"/houses/:houseID/modules/:moduleID/turn-on",
			func(c *gin.Context) { presentation.TurnOnModule(c, appContainer) },
		)
		moduleGroup.PUT(
			"/houses/:houseID/modules/:moduleID/turn-off",
			func(c *gin.Context) { presentation.TurnOffModule(c, appContainer) },
		)
	}

	go func() {
		for {
			event, err := appContainer.ModuleService.GetModuleVerificationEvent(ctx)
			if err != nil {
				log.Printf("Error while reading message: %v", err)
				continue
			}

			log.Printf("Event successfully processed: %v", event)
		}
	}()

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

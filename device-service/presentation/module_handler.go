package presentation

import (
	"device-service/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAvailableModules(c *gin.Context, container *shared.Container) {
	modules, err := container.ModuleService.GetAllModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve modules"})
		return
	}

	c.JSON(http.StatusOK, modules)
}

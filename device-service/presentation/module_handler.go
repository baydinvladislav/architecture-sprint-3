package presentation

import (
	web_schemas "device-service/presentation/web-schemas"
	"device-service/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func GetModulesByHouseId(c *gin.Context, container *shared.Container) {
	houseIDStr := c.Param("houseID")

	houseID, err := uuid.Parse(houseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	modules, err := container.ModuleService.GetModulesByHouseID(houseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve modules"})
		return
	}

	c.JSON(http.StatusOK, modules)
}

func AddModuleToHouse(c *gin.Context, container *shared.Container) {
	houseIDStr := c.Param("houseID")
	houseID, err := uuid.Parse(houseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	var newModule web_schemas.ConnectModuleIn

	if err := c.ShouldBindJSON(&newModule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newModuleResponse, err := container.ModuleService.AddModuleToHouse(houseID, newModule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add module"})
		return
	}

	c.JSON(http.StatusCreated, newModuleResponse)
}

func TurnOnModule(c *gin.Context, container *shared.Container) {
	houseIDStr := c.Param("houseID")
	moduleIDStr := c.Param("moduleID")

	houseID, err := uuid.Parse(houseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	if err := container.ModuleService.TurnOnModule(houseID, moduleID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found or could not be turned on"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module turned on successfully"})
}

func TurnOffModule(c *gin.Context, container *shared.Container) {
	houseIDStr := c.Param("houseID")
	moduleIDStr := c.Param("moduleID")

	houseID, err := uuid.Parse(houseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	if err := container.ModuleService.TurnOffModule(houseID, moduleID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found or could not be turned off"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module turned off successfully"})
}

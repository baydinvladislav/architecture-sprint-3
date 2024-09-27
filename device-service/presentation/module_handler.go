package presentation

import (
	web_schemas "device-service/presentation/web-schemas"
	"device-service/shared"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	houseID := c.Param("houseID")

	modules, err := container.ModuleService.GetModulesByHouseID(houseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve modules"})
		return
	}

	c.JSON(http.StatusOK, modules)
}

func TurnOnModule(c *gin.Context, container *shared.Container) {
	houseID := c.Param("houseID")
	moduleID := c.Param("moduleID")

	if err := container.ModuleService.TurnOnModule(houseID, moduleID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found or could not be turned on"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module turned on successfully"})
}

func TurnOffModule(c *gin.Context, container *shared.Container) {
	houseIDStr := c.Param("houseID")
	moduleIDStr := c.Param("moduleID")

	houseID, err := strconv.ParseUint(houseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	if err := container.ModuleService.TurnOffModule(uint(houseID), uint(moduleID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found or could not be turned off"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module turned off successfully"})
}

func AddModuleToHouse(c *gin.Context, container *shared.Container) {
	houseID := c.Param("houseID")
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

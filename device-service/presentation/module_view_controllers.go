package presentation

import (
	"device-service/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// GetAvailableModules godoc
// @Summary Получение доступных модулей
// @Description Получить список всех доступных модулей
// @Tags modules
// @Produce json
// @Success 200 {array} web_schemas.ModuleOut
// @Failure 500 {object} ErrorResponse
// @Router /modules [get]
func GetAvailableModules(c *gin.Context, container *shared.Container) {
	modules, err := container.ModuleService.GetAllModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve modules"})
		return
	}

	c.JSON(http.StatusOK, modules)
}

// GetModulesByHouseId godoc
// @Summary Получение модулей по ID дома
// @Description Получить все модули, привязанные к определенному дому
// @Tags modules
// @Produce json
// @Param houseID path string true "House ID"
// @Success 200 {array} web_schemas.ModuleOut
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /houses/{houseID}/modules [get]
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

func GetModuleState(c *gin.Context, container *shared.Container) {
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

	stateInfo, err := container.ModuleService.GetModuleState(houseID, moduleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not fetch state due to an error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": stateInfo})
}

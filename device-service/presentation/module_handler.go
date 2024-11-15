package presentation

import (
	"context"
	"device-service/repository"
	"device-service/service"
	"device-service/shared"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
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

// AddModuleToHouse godoc
// @Summary Добавление модуля в дом
// @Description Запрос на добавление модуля в дом
// @Tags modules
// @Produce json
// @Param houseID path string true "House ID"
// @Param moduleID path string true "Module ID"
// @Success 201 {object} web_schemas.ModuleOut
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /houses/{houseID}/modules/{moduleID} [post]
func AddModuleToHouse(c *gin.Context, container *shared.Container) {
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

	newModuleResponse, err := container.ModuleService.RequestAdditionModuleToHouse(houseID, moduleID)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add module"})
		return
	}

	c.JSON(http.StatusCreated, newModuleResponse)
}

// TurnOnModule godoc
// @Summary Включение модуля
// @Description Включить модуль, привязанный к определенному дому
// @Tags modules
// @Produce json
// @Param houseID path string true "House ID"
// @Param moduleID path string true "Module ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /houses/{houseID}/modules/{moduleID}/turn-on [post]
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
		if errors.Is(err, repository.ErrModuleAlreadyOn) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Module is already turned on"})
			return
		} else if errors.Is(err, repository.ErrModuleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
			return
		} else if errors.Is(err, service.ErrKafkaSupplier) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Kafka message"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not turn on module"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module turned on successfully"})
}

// TurnOffModule godoc
// @Summary Выключение модуля
// @Description Выключить модуль, привязанный к определенному дому
// @Tags modules
// @Produce json
// @Param houseID path string true "House ID"
// @Param moduleID path string true "Module ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /houses/{houseID}/modules/{moduleID}/turn-off [post]
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
		if errors.Is(err, repository.ErrModuleAlreadyOff) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Module is already turned off"})
			return
		} else if errors.Is(err, repository.ErrModuleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
			return
		} else if errors.Is(err, service.ErrKafkaSupplier) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Kafka message"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not turn off module"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module turned off successfully"})
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

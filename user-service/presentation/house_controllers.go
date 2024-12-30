package presentation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	web_schemas "user-service/schemas/web"
	"user-service/shared"
)

// CreateUserHouse godoc
// @Summary Создание нового дома для пользователя
// @Description Создать новый дом, связанный с пользователем, и вернуть информацию о доме
// @Tags houses
// @Accept json
// @Produce json
// @Param house body web_schemas.NewHouseIn true "New House"
// @Success 201 {object} web_schemas.HouseOut
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /houses [post]
func CreateUserHouse(c *gin.Context, container *shared.Container) {
	var house web_schemas.NewHouseIn

	if err := c.ShouldBindJSON(&house); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not get username from JWT token"})
		return
	}

	user, err := container.UserService.GetByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not exists"})
		return
	}

	newHouse, err := container.HouseService.CreateUserHouse(user.ID, house)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create house"})
		return
	}

	c.JSON(http.StatusCreated, newHouse)
}

// GetUserHouses godoc
// @Summary Получение всех домов пользователя
// @Description Получить список всех домов, связанных с пользователем
// @Tags houses
// @Produce json
// @Success 200 {array} web_schemas.HouseOut
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /houses [get]
func GetUserHouses(c *gin.Context, container *shared.Container) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not get username from JWT token"})
		return
	}

	user, err := container.UserService.GetByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	houses, err := container.HouseService.GetUserHouses(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user's houses"})
		return
	}

	if len(houses) == 0 {
		c.JSON(http.StatusOK, []web_schemas.HouseOut{})
		return
	}

	c.JSON(http.StatusOK, houses)
}

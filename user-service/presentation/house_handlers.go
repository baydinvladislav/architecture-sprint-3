package presentation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	web_schemas "user-service/presentation/web-schemas"
	"user-service/shared"
)

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

	newHouseResponse := web_schemas.HouseOut{
		ID:      newHouse.ID,
		Address: newHouse.Address,
		Square:  newHouse.Square,
		UserID:  newHouse.UserID,
	}

	c.JSON(http.StatusCreated, newHouseResponse)
}

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

func UpdateUserHouse(c *gin.Context, container *shared.Container) {}

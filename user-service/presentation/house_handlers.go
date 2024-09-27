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

	newHouse, err := container.HouseService.CreateUserHouse(house)
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

func GetHouseById(c *gin.Context, container *shared.Container) {
	houseId := c.Param("Id")

	house, err := container.HouseService.GetUserHouse(houseId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "house not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"house": house})
}

func UpdateUserHouse(c *gin.Context, container *shared.Container) {

}

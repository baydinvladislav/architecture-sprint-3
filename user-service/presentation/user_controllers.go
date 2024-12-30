package presentation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/shared"
)

// GetUserById godoc
// @Summary Получение пользователя по ID
// @Description Получить данные пользователя по его ID
// @Tags users
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{} "user"
// @Failure 404 {object} ErrorResponse
// @Router /users/{userId} [get]
func GetUserById(c *gin.Context, container *shared.Container) {
	userId := c.Param("userId")

	user, err := container.UserService.GetCurrent(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

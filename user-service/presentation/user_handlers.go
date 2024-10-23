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

// GetUserByUsername godoc
// @Summary Получение пользователя по имени пользователя
// @Description Получить данные пользователя по его имени пользователя (username)
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} map[string]interface{} "user"
// @Failure 404 {object} ErrorResponse
// @Router /users/username/{username} [get]
func GetUserByUsername(c *gin.Context, container *shared.Container) {
	// пока не пригодился, пока заглушка
}

// UpdateUser godoc
// @Summary Обновление данных пользователя
// @Description Обновить данные существующего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body web_schemas.LoginRequest true "Updated User Data"
// @Success 200 {object} map[string]interface{} "user"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{userId} [put]
func UpdateUser(c *gin.Context, container *shared.Container) {
	// пока не пригодился, пока заглушка
}

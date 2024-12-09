package presentation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/schemas/web"
	"user-service/shared"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Зарегистрировать нового пользователя и получить Access и Refresh токены
// @Tags users
// @Accept json
// @Produce json
// @Param user body web_schemas.UserIn true "New User"
// @Success 201 {object} web_schemas.UserClaimsOut
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context, container *shared.Container) {
	var userData web.UserIn

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newUser, err := container.UserService.SignUp(&userData)
	//здесь должен быть свитч на разные возможные ошибки
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// LoginUser godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя с генерацией Access и Refresh токенов
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body web_schemas.LoginRequest true "Login credentials"
// @Success 200 {object} web_schemas.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func LoginUser(c *gin.Context, container *shared.Container) {
	var credentials web.LoginRequest

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authUser, err := container.UserService.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error during user login"})
		return
	}

	c.JSON(http.StatusOK, authUser)
}

// RefreshToken godoc
// @Summary Обновление токена доступа
// @Description Обновить Access и Refresh токены с использованием действующего Refresh токена
// @Tags users
// @Accept json
// @Produce json
// @Param user body web_schemas.RefreshTokenRequest true "New User"
// @Success 200 {object} map[string]string "Access и Refresh токены"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /refresh-token [post]
func RefreshToken(c *gin.Context, container *shared.Container) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newAccessToken, newRefreshToken, err := container.UserService.UpdateUserRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error during tokens refreshing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

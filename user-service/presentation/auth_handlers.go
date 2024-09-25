package presentation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/presentation/web-schemas"
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
// @Param user body web_schemas.NewUserIn true "New User"
// @Success 201 {object} web_schemas.NewUserOut
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /users/register [post]
func RegisterUser(c *gin.Context, container *shared.Container) {
	var user web_schemas.NewUserIn

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := container.UserService.SignUp(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	accessToken, err := container.AuthService.GenerateAccessToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := container.AuthService.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	newUserResponse := web_schemas.NewUserOut{
		// TODO: сходить в БД и получить ID, а лучше научиться в Go писать в таблицу и получать результат записи
		//ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusCreated, newUserResponse)
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
// @Router /users/login [post]
func LoginUser(c *gin.Context, container *shared.Container) {
	var credentials web_schemas.LoginRequest

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := container.UserService.GetByUsername(credentials.Username)
	if err != nil || user.Password != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := container.AuthService.GenerateAccessToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := container.AuthService.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	response := web_schemas.LoginResponse{
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context, container *shared.Container) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	claims, err := container.AuthService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	accessToken, err := container.AuthService.GenerateAccessToken(claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := container.AuthService.GenerateRefreshToken(claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

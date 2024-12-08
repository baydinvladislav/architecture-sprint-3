package web

import "github.com/google/uuid"

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

package dto

import "github.com/google/uuid"

type UserDtoSchema struct {
	ID       uuid.UUID
	Username string
	Password string
	Email    string
}

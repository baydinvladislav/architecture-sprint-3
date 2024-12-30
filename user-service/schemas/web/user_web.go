package web

import "github.com/google/uuid"

type UserIn struct {
	Username string
	Password string
}

type UserOut struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

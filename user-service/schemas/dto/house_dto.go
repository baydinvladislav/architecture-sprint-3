package dto

import "github.com/google/uuid"

type HouseDtoSchema struct {
	ID      uuid.UUID
	Address string
	Square  float64
	UserID  uuid.UUID
}

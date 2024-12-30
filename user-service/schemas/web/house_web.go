package web

import "github.com/google/uuid"

type NewHouseIn struct {
	Address string  `json:"address"`
	Square  float64 `json:"square"`
}

type UpdateHouseIn struct {
	HouseID uuid.UUID `json:"house_id"`
	Address string    `json:"address"`
	Square  float64   `json:"square"`
	UserID  uuid.UUID `json:"user_id"`
}

type HouseOut struct {
	ID      uuid.UUID `json:"house_id"`
	Address string    `json:"address"`
	Square  float64   `json:"square"`
	UserID  uuid.UUID `json:"user_id"`
}

package web_schemas

type NewHouseIn struct {
	Address string  `json:"address"`
	Square  float64 `json:"square"`
	UserID  uint    `json:"user_id"`
}

type UpdateHouseIn struct {
	HouseID uint    `json:"house_id"`
	Address string  `json:"address"`
	Square  float64 `json:"square"`
	UserID  uint    `json:"user_id"`
}

type HouseOut struct {
	ID      uint    `json:"house_id"`
	Address string  `json:"address"`
	Square  float64 `json:"square"`
	UserID  uint    `json:"user_id"`
}

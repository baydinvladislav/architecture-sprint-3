package web_schemas

import "time"

type ModuleOut struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

type ConnectModuleIn struct {
	ID uint `json:"id"`
}

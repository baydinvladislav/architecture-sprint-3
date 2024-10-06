package schemas

type Event struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type EventPayload interface {
	IsEventPayload()
}

type ModuleAdditionPayload struct {
	HouseID  string `json:"source_id"`
	ModuleID string `json:"source_type"`
	UserID   string `json:"value"`
	Time     int64  `json:"time"`
}

func (ModuleAdditionPayload) IsEventPayload() {}

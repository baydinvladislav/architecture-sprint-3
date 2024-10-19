package schemas

type Event struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type EventPayload interface {
	IsEventPayload()
}

type ModuleVerifyPayload struct {
	Time     int64  `json:"time"`
	HomeID   string `json:"home_id"`
	ModuleID string `json:"module_id"`
}

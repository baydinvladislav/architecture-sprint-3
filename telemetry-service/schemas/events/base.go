package events

type Event struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type BasePayload interface {
	IsEventPayload()
}

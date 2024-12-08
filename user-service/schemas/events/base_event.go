package events

type BaseEvent struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type EventPayload interface {
	IsEventPayload()
}

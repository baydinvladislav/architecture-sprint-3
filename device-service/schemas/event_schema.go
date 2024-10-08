package schemas

type Event struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type EventPayload interface {
	IsEventPayload()
}

type ModuleVerification struct {
	HouseID  string `json:"source_id"`
	ModuleID string `json:"source_type"`
	UserID   string `json:"value"`
	Time     int64  `json:"time"`
	Decision string `json:"decision"`
}

func (ModuleVerification) IsEventPayload() {}

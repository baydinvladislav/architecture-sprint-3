package schemas

type BaseEvent struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type EventPayload interface {
	IsEventPayload()
}

type HomeVerificationEvent struct {
	HouseID  string `json:"source_id"`
	ModuleID string `json:"source_type"`
	Time     int64  `json:"time"`
}

func (HomeVerificationEvent) IsEventPayload() {}

type ModuleVerificationEvent struct {
	HouseID  string `json:"source_id"`
	ModuleID string `json:"source_type"`
	UserID   string `json:"value"`
	Time     int64  `json:"time"`
	Decision string `json:"decision"`
}

func (ModuleVerificationEvent) IsEventPayload() {}

type ChangeEquipmentStateEvent struct {
	HouseID  string                 `json:"house_id"`
	ModuleID string                 `json:"module_id"`
	Time     int64                  `json:"time"`
	State    map[string]interface{} `json:"state"`
}

func (ChangeEquipmentStateEvent) IsEventPayload() {}

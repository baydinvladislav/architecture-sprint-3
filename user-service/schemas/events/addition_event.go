package events

type ModuleAdditionEvent struct {
	HouseID  string `json:"source_id"`
	ModuleID string `json:"source_type"`
	Time     int64  `json:"time"`
}

func (ModuleAdditionEvent) IsEventPayload() {}

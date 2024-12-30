package events

type ModuleVerificationEvent struct {
	HouseID  string `json:"source_id"`
	ModuleID string `json:"source_type"`
	UserID   string `json:"value"`
	Time     int64  `json:"time"`
	Decision string `json:"decision"`
}

func (ModuleVerificationEvent) IsEventPayload() {}

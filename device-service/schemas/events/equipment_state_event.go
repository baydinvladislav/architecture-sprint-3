package events

type ChangeEquipmentStateEvent struct {
	HouseID  string                 `json:"house_id"`
	ModuleID string                 `json:"module_id"`
	Time     int64                  `json:"time"`
	State    map[string]interface{} `json:"state"`
}

func (ChangeEquipmentStateEvent) IsEventPayload() {}

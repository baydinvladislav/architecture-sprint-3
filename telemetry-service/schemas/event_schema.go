package schemas

type Event struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

const (
	SourceTypeSensor    = "sensor"
	SourceTypeEquipment = "equipment"
)

type EventPayload interface {
	IsEventPayload()
}

type TelemetryPayload struct {
	SourceID   string  `json:"source_id"`
	SourceType string  `json:"source_type"`
	Value      float64 `json:"value"`
	Time       int64   `json:"time"`
}
type EmergencyPayload struct {
	EquipmentID string `json:"equipment_id"`
	Reason      string `json:"reason"`
}

func (TelemetryPayload) IsEventPayload() {}
func (EmergencyPayload) IsEventPayload() {}

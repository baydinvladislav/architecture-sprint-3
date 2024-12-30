package events

const (
	SourceTypeSensor    = "sensor"
	SourceTypeEquipment = "equipment"
)

// TelemetryPayload общий ивент с телеметрией
type TelemetryPayload struct {
	SourceID   string  `json:"source_id"`
	SourceType string  `json:"source_type"`
	Value      float64 `json:"value"`
	Time       int64   `json:"time"`
}

func (TelemetryPayload) IsEventPayload() {}

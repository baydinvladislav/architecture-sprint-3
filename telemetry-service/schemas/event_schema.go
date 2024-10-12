package schemas

import "github.com/google/uuid"

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

// TelemetryPayload общий ивент с телеметрией
type TelemetryPayload struct {
	SourceID   string  `json:"source_id"`
	SourceType string  `json:"source_type"`
	Value      float64 `json:"value"`
	Time       int64   `json:"time"`
}

func (TelemetryPayload) IsEventPayload() {}

// InstallModuleToHousePayload фиксируем ивентом иницилизации только что подключенный дом к системе SmartHome,
// в основном для реализации SAGA
type InstallModuleToHousePayload struct {
	HouseID  uint      `json:"house_id"`
	ModuleID uuid.UUID `json:"module_id"`
}

func (InstallModuleToHousePayload) IsEventPayload() {}

// EmergencyPayload предусмотрим обработку экстренного выключения оборудования, поступающие в критические моменты,
// например температура в спальне достигла 30 градусов -> выключим модуль отопление,
// дым в гостивной -> отключить всё электро-оборудование в гостиной и т.д.
// в основном для реализации межсерисного взаимодействия
type EmergencyPayload struct {
	SourceID string `json:"source_id"` // Идентификатор устройства, отправившего экстренное сообщение
	Reason   string `json:"reason"`    // Причина срабатывания экстренной ситуации
}

func (EmergencyPayload) IsEventPayload() {}

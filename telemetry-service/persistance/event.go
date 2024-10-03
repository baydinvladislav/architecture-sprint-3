package persistance

type Event struct {
	EquipmentID string `json:"equipment_id"`
	Reason      string `json:"reason"`
	Time        int64  `json:"time"`
}

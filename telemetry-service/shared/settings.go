package shared

type AppSettings struct {
	KafkaBroker      string
	TelemetryTopic   string
	EmergencyTopic   string
	GroupID          string
	MongoURI         string
	DatabaseName     string
	CollectionName   string
	DeviceServiceUrl string
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		KafkaBroker:      "kafka:9092",
		TelemetryTopic:   "telemetry.data",
		EmergencyTopic:   "forced.module.shutdown",
		GroupID:          "telemetry_group",
		MongoURI:         "mongodb://root:mongodb@mongo:27017",
		DatabaseName:     "telemetry_database",
		CollectionName:   "events",
		DeviceServiceUrl: "http://device-service:8081",
	}
}

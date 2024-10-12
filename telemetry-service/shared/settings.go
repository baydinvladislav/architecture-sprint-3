package shared

type AppSettings struct {
	// broker
	KafkaBroker    string
	EmergencyTopic string
	NewHouseTopic  string
	TelemetryTopic string
	GroupID        string

	// db
	MongoURI            string
	DatabaseName        string
	TelemetryCollection string
	HouseCollection     string

	// external
	DeviceServiceUrl string
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		KafkaBroker:    "kafka:9092",
		TelemetryTopic: "telemetry.data.topic",
		EmergencyTopic: "forced.module.shutdown.topic",
		NewHouseTopic:  "house.initialization.topic",
		GroupID:        "telemetry_group",

		MongoURI:            "mongodb://root:mongodb@mongo:27017",
		DatabaseName:        "telemetry_database",
		TelemetryCollection: "events",

		DeviceServiceUrl: "http://device-service:8081",
	}
}

package shared

import (
	"os"
)

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
		KafkaBroker:    getEnv("KAFKA_BROKER", "kafka:9092"),
		TelemetryTopic: getEnv("TELEMETRY_TOPIC", "telemetry.data"),
		EmergencyTopic: getEnv("EMERGENCY_TOPIC", "forced.module.shutdown"),
		NewHouseTopic:  getEnv("NEW_HOUSE_TOPIC", "house.initialization"),
		GroupID:        getEnv("KAFKA_GROUP_ID", "telemetry_group"),

		MongoURI:            getEnv("MONGO_URI", "mongodb://root:mongodb@mongo:27017/telemetry_database?authSource=admin"),
		DatabaseName:        getEnv("MONGO_DATABASE_NAME", "telemetry_database"),
		TelemetryCollection: getEnv("MONGO_TELEMETRY_COLLECTION", "events"),

		DeviceServiceUrl: getEnv("DEVICE_SERVICE_URL", "http://device-service:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

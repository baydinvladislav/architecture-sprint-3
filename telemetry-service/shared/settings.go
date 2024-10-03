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
		KafkaBroker:      "",
		MongoURI:         "",
		DatabaseName:     "",
		CollectionName:   "",
		DeviceServiceUrl: "",
	}
}

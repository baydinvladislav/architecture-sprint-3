package shared

type AppSettings struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string

	KafkaBroker                  string
	ModuleAddedKafkaTopic        string
	ModuleVerificationKafkaTopic string
	KafkaGroupID                 string
}

func (s *AppSettings) DSN() string {
	return "host=" + s.DBHost +
		" user=" + s.DBUser +
		" password=" + s.DBPassword +
		" dbname=" + s.DBName +
		" port=" + s.DBPort +
		" sslmode=" + s.DBSSLMode
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		DBHost:     "db",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "sprint_3",
		DBPort:     "5432",
		DBSSLMode:  "disable",

		KafkaBroker:                  "kafka:9092",
		ModuleAddedKafkaTopic:        "module.addition.topic",
		ModuleVerificationKafkaTopic: "module.verification.topic",
		KafkaGroupID:                 "user_service_group",
	}
}

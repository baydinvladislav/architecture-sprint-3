package shared

type AppSettings struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string

	KafkaBrokers            []string
	ModuleAdditionTopic     string
	ModuleVerificationTopic string
	KafkaGroupID            string
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
		DBHost:     "user-db",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "user_db",
		DBPort:     "5432",
		DBSSLMode:  "disable",

		KafkaBrokers:            []string{"kafka1:9092", "kafka2:9093", "kafka3:9094"},
		ModuleAdditionTopic:     "module.addition.topic",
		ModuleVerificationTopic: "module.verification.topic",
		KafkaGroupID:            "user_service_group",
	}
}

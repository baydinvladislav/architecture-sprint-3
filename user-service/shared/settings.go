package shared

type AppSettings struct {
	MinHomeSquare float64
	MaxHomeSquare float64

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string

	KafkaBroker             string
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
		MinHomeSquare: 20.0,
		MaxHomeSquare: 200.0,

		DBHost:     "db",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "sprint_3",
		DBPort:     "5432",
		DBSSLMode:  "disable",

		KafkaBroker:             "kafka:9092",
		ModuleAdditionTopic:     "module.addition.topic",
		ModuleVerificationTopic: "module.verification.topic",
		KafkaGroupID:            "user_service_group",
	}
}

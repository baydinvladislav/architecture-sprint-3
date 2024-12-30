package repository

type TelemetryRepositoryInterface interface {
	InsertEvent(event interface{}) error
	Close()
}

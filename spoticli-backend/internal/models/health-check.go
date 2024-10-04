package models

// TODO: move models to their own go module
type HealthCheck struct {
	ServiceName string
	ManagedBy   string
	Timestamp   string
	Version     string
}

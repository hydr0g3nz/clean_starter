package infra

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{}) // Fatal should typically cause process exit

	// Creates a child logger with added fields
	With(fields map[string]interface{}) Logger // Return the interface type

	// Sync flushes any buffered log entries.
	Sync() error

	// Close() error // Optional, if needed to close file handles etc.
}

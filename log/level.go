package log

// LogLevel specifies the severity of a given log message
type LogLevel string

// Log levels
const (
	LogLevelDebug   LogLevel = "debug"
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warn"
	LogLevelError   LogLevel = "error"
)

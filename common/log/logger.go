package log

type LogLevel uint8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	NONE
)

type Logger interface {
	// Info - Log info message
	Info(v ...interface{})

	// Debug - Log Debug message
	Debug(v ...interface{})

	// Warning - Log Warning message
	Warning(v ...interface{})

	// Trace - Log Trace message
	Trace(v ...interface{})

	// Error - Log Error message
	Error(v ...interface{})

	// ErrorTrace - Log error with trace
	ErrorTrace(e error)

	// Panic - Log error with trace
	Panic(v ...interface{})
}

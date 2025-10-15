package logger


import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
)

// Logger interface defines logging methods
type Logger interface {
	Debug(ctx context.Context, message string, fields map[string]interface{})
	Info(ctx context.Context, message string, fields map[string]interface{})
	Warn(ctx context.Context, message string, fields map[string]interface{})
	Error(ctx context.Context, message string, err error, fields map[string]interface{})
	Fatal(ctx context.Context, message string, err error, fields map[string]interface{})
}

// LogLevel represents logging levels
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// StandardLogger implements Logger interface using Go's standard log package
type StandardLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewStandardLogger creates a new standard logger
func NewStandardLogger(level LogLevel) Logger {
	return &StandardLogger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// logEntry represents a structured log entry
type logEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
}

func (l *StandardLogger) log(ctx context.Context, level LogLevel, levelName, message string, err error, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := logEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     levelName,
		Message:   message,
		Fields:    fields,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// Extract trace ID from context if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if str, ok := traceID.(string); ok {
			entry.TraceID = str
		}
	}

	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf("Error marshaling log entry: %v", err)
		return
	}

	l.logger.Println(string(data))
}

func (l *StandardLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, DebugLevel, "DEBUG", message, nil, fields)
}

func (l *StandardLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, InfoLevel, "INFO", message, nil, fields)
}

func (l *StandardLogger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, WarnLevel, "WARN", message, nil, fields)
}

func (l *StandardLogger) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	l.log(ctx, ErrorLevel, "ERROR", message, err, fields)
}

func (l *StandardLogger) Fatal(ctx context.Context, message string, err error, fields map[string]interface{}) {
	l.log(ctx, FatalLevel, "FATAL", message, err, fields)
	os.Exit(1)
}

// ParseLogLevel parses log level from string
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug", "DEBUG":
		return DebugLevel
	case "info", "INFO":
		return InfoLevel
	case "warn", "WARN", "warning", "WARNING":
		return WarnLevel
	case "error", "ERROR":
		return ErrorLevel
	case "fatal", "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// DefaultLogger creates a default logger with INFO level
func DefaultLogger() Logger {
	return NewStandardLogger(InfoLevel)
}
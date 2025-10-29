package model

import "time"

type LogLevel string

func (l LogLevel) String() string {
	return string(l)
}

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`  // log timestamp
	Level     LogLevel          `json:"level"`      // e.g., INFO, WARN, ERROR
	Component string            `json:"component"`  // e.g., cache, api, db
	Host      string            `json:"host"`       // e.g., web01, db01
	RequestID string            `json:"request_id"` // e.g., req-abc123
	Message   string            `json:"message"`    // log message
	Metadata  map[string]string `json:"metadata"`   // Dynamic fields
	Log       string            `json:"log"`        // Raw log
}

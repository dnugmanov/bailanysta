package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level  Level
	writer io.Writer
}

type LogEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Fields  Fields    `json:"fields,omitempty"`
	File    string    `json:"file,omitempty"`
	Line    int       `json:"line,omitempty"`
}

type Fields map[string]interface{}

func New(level string, writer io.Writer) *Logger {
	if writer == nil {
		writer = os.Stdout
	}

	var lvl Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = DebugLevel
	case "info":
		lvl = InfoLevel
	case "warn":
		lvl = WarnLevel
	case "error":
		lvl = ErrorLevel
	case "fatal":
		lvl = FatalLevel
	default:
		lvl = InfoLevel
	}

	return &Logger{
		level:  lvl,
		writer: writer,
	}
}

func (l *Logger) log(level Level, message string, fields Fields) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Time:    time.Now(),
		Level:   level.String(),
		Message: message,
		Fields:  fields,
	}

	// Add caller info for non-Info levels
	if level != InfoLevel {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry.File = file
			entry.Line = line
		}
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	fmt.Fprintln(l.writer, string(jsonData))
}

func (l *Logger) Debug(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DebugLevel, message, f)
}

func (l *Logger) Info(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(InfoLevel, message, f)
}

func (l *Logger) Warn(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WarnLevel, message, f)
}

func (l *Logger) Error(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ErrorLevel, message, f)
}

func (l *Logger) Fatal(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FatalLevel, message, f)
	os.Exit(1)
}

// Middleware compatible logger
func (l *Logger) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.Info(message)
}

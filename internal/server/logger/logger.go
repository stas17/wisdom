package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// LogLevel defines on what level logs should be recorded.
type LogLevel int

const (
	// DebugLevel defines debug log level.
	DebugLevel LogLevel = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// DisabledLevel disables the logger
	DisabledLevel
)

var levels = map[string]LogLevel{
	"debug":    DebugLevel,
	"info":     InfoLevel,
	"warn":     WarnLevel,
	"error":    ErrorLevel,
	"fatal":    FatalLevel,
	"panic":    PanicLevel,
	"no_level": NoLevel,
	"disabled": DisabledLevel,
}

// LevelFromString returns a LogLevel according to string value from config.
func LevelFromString(s string) LogLevel {
	if logLevel, ok := levels[s]; ok {
		return logLevel
	}

	// use info level as default.
	return InfoLevel
}

// Logger a struct that wraps zerolog.
type Logger struct {
	*zerolog.Logger
}

// New sets up and returns a new Logger wrapper around a zerolog logger.
func New(level LogLevel) *Logger {
	var l zerolog.Level
	switch level {
	case DebugLevel:
		l = zerolog.DebugLevel
	case InfoLevel:
		l = zerolog.InfoLevel
	case FatalLevel:
		l = zerolog.FatalLevel
	case DisabledLevel:
		l = zerolog.Disabled
	case WarnLevel:
		l = zerolog.WarnLevel
	case ErrorLevel:
		l = zerolog.ErrorLevel
	case PanicLevel:
		l = zerolog.PanicLevel
	case NoLevel:
		l = zerolog.NoLevel
	default:
		// This default value may be changed depending on the requirements of the service.
		l = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).Level(l).With().Timestamp().Logger()
	return &Logger{&logger}
}

// WrapHook wraps the Hook function from zerolog.
func (l Logger) WrapHook(h zerolog.Hook) Logger {
	hooked := l.Hook(h)
	return Logger{&hooked}
}

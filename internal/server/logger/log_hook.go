package logger

import (
	"github.com/rs/zerolog"
)

// SourceLogHook represents a hook for the logger package to add custom fields.
type SourceLogHook struct {
	Source string
}

// Run runs the hook with the event about logs source.
func (h SourceLogHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if level != zerolog.NoLevel {
		e.Str("source", h.Source)
	}
}

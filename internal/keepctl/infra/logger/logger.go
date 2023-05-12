// Package logger configures logging facility.
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger represents logger entity.
type Logger struct {
	zerolog.Logger
}

// New creates new logger entity.
func New(debug bool) *Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{
		Out:             os.Stdout,
		FormatTimestamp: func(i interface{}) string { return "" },
		FormatLevel:     func(i interface{}) string { return "" },
	}

	l := zerolog.New(output).
		With().
		Logger()

	return &Logger{l}
}

// Package logger configures logging facility.
package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger represents logger entity.
type Logger struct {
	zerolog.Logger
}

// New creates new logger entity.
func New(level string) (*Logger, error) {
	if level == "" {
		level = zerolog.LevelInfoValue
	}

	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("logger - New - zerolog.ParseLevel: %w", err)
	}

	zerolog.SetGlobalLevel(lv)

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.TimeFormat = time.RFC822

	skipFrameCount := 3
	l := zerolog.New(output).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{l}, nil
}

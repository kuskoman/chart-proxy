package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// ErrUnknownLogFormat is returned when the log format is not recognized by the slog package
var ErrUnknownLogFormat = fmt.Errorf("unknown log format")

func setupSlog(loggingConfig *LoggingConfig) (*slog.Logger, error) {
	level := slog.LevelInfo
	err := level.UnmarshalText([]byte(loggingConfig.Level))
	if err != nil {
		return nil, err
	}

	var handler slog.Handler
	switch strings.ToLower(loggingConfig.Format) {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	default:
		return nil, ErrUnknownLogFormat
	}

	return slog.New(handler), nil
}

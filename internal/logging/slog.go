package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/kuskoman/chart-proxy/pkg/config"
)

// ErrUnknownLogFormat is returned when the log format is not recognized by the slog package
var ErrUnknownLogFormat = fmt.Errorf("unknown log format")

// SetupSlog configures the slog package with the given configuration.
// The function is designed to work as a reload hook for the config manager.
func SetupSlog(config *config.Config) error {
	loggingConfig := config.Logging

	level := slog.LevelInfo
	err := level.UnmarshalText([]byte(loggingConfig.Level))
	if err != nil {
		return err
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
		return ErrUnknownLogFormat
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

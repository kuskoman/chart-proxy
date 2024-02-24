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

// SetupLogin configures the logging for the application based on the configuration.
// Can be called multiple times to reload the logging configuration.
func SetupLogging(cfg *config.Config, errChan chan error) {
	err := setupSlog(&cfg.Logging)
	if err != nil {
		errChan <- err
	}
}

func setupSlog(loggingConfig *config.LoggingConfig) error {
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

	slog.Debug("logging configured", "level", level, "format", loggingConfig.Format)

	return nil
}

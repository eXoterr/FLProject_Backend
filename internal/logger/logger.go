package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

const (
	TextFormat = "text"
	JSONFormat = "json"

	LevelDebug string = "debug"
	LevelInfo  string = "info"
	LevelWarn  string = "warn"
	LevelError string = "error"
)

func SetupLogger(level string, format string) (*slog.Logger, error) {
	var handler slog.Handler

	logLevel, err := parseLogLevel(level)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	switch format {
	case TextFormat:
		handler = slog.NewTextHandler(os.Stdout, opts)
	case JSONFormat:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		return nil, fmt.Errorf("unknown logger format: %s", format)
	}

	logger := slog.New(handler)

	return logger, nil
}

func MustSetupLogger(level string, format string) *slog.Logger {
	logger, err := SetupLogger(level, format)
	if err != nil {
		log.Fatalf("unable to setup logger: %s", err)
	}
	return logger
}

func parseLogLevel(level string) (slog.Level, error) {
	switch level {
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

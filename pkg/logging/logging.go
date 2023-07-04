package logging

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type LoggerDeps struct {
	LogLevel string
	LogFile  string
	LogSize  int
	LogAge   int
}

func NewLogger(cfg *LoggerDeps) (*zerolog.Logger, error) {
	var writer io.Writer = os.Stderr
	if cfg.LogFile != "" {
		writer = &lumberjack.Logger{
			Filename:   cfg.LogFile,
			MaxSize:    cfg.LogSize,
			MaxBackups: 0,
			MaxAge:     cfg.LogAge,
		}
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, errors.Wrap(err, "parse log level")
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger, nil
}

func NewDefaultLogger() zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return logger
}

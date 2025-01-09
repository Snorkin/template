package logger

import (
	"example/config"
	"github.com/rs/zerolog"
	"os"
)

func InitLogger() {
	logger := initLogger()
	Log = &appLogger{logger: logger}
}

func initLogger() *zerolog.Logger {
	cfg := config.GetConfig()

	var w zerolog.LevelWriter
	if cfg.Environment == "dev" {
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		w = zerolog.MultiLevelWriter(os.Stdout)
	}

	logger := zerolog.New(w).Level(loggerLevelMap[cfg.Logger.Level]).With().
		CallerWithSkipFrameCount(cfg.Logger.SkipFrameCount).Timestamp().Logger()
	logger.Info().Msgf("Logger initialized: level %s", cfg.Logger.Level)
	return &logger
}

package logger

import (
	"example/config"
	"github.com/rs/zerolog"
	"os"
)

var loggerLevelMap = map[string]zerolog.Level{
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"panic":    zerolog.PanicLevel,
	"fatal":    zerolog.FatalLevel,
	"noLevel":  zerolog.NoLevel,
	"disabled": zerolog.Disabled,
}

type appLogger struct {
	cfg    config.Config
	logger *zerolog.Logger
}

func NewLogger(cfg config.Config) Logger {
	return &appLogger{cfg: cfg, logger: initLogger(cfg)}
}

// TODO: add env check
func initLogger(cfg config.Config) *zerolog.Logger {
	w := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	logger := zerolog.New(w).Level(loggerLevelMap[cfg.Logger.Level]).With().
		CallerWithSkipFrameCount(cfg.Logger.SkipFrameCount).Timestamp().Logger()
	return &logger
}

func (a *appLogger) Debug(msg string) {
	a.logger.Debug().Msg(msg)
}

func (a *appLogger) Debugf(template string, args ...interface{}) {
	a.logger.Debug().Msgf(template, args...)
}

func (a *appLogger) Info(msg string) {
	a.logger.Info().Msg(msg)
}

func (a *appLogger) Infof(template string, args ...interface{}) {
	a.logger.Info().Msgf(template, args...)
}

func (a *appLogger) Warn(msg string) {
	a.logger.Warn().Msg(msg)
}

func (a *appLogger) Warnf(template string, args ...interface{}) {
	a.logger.Warn().Msgf(template, args...)
}

func (a *appLogger) Error(err error) {
	a.logger.Error().Msg(err.Error())
}

func (a *appLogger) Errorf(template string, args ...interface{}) {
	a.logger.Error().Msgf(template, args...)
}

func (a *appLogger) Panic(msg string) {
	a.logger.Panic().Msg(msg)
}

func (a *appLogger) Panicf(template string, args ...interface{}) {
	a.logger.Panic().Msgf(template, args...)
}

func (a *appLogger) Fatal(msg string) {
	a.logger.Fatal().Msg(msg)
}

func (a *appLogger) Fatalf(template string, args ...interface{}) {
	a.logger.Fatal().Msgf(template, args...)
}

package logger

import (
	"example/config"
	"github.com/rs/zerolog"
	"os"
	"reflect"
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

var Log *appLogger

type appLogger struct {
	logger *zerolog.Logger
}

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
	logger.Info().Msgf("Logger initialized: level - %s", cfg.Logger.Level)
	return &logger
}

func (a *appLogger) Debug(msg string) {
	a.logger.Debug().Msg(msg)
}

func (a *appLogger) Debugf(msg string, keyValue ...interface{}) {
	if len(keyValue)%2 != 0 {
		a.logger.Error().Msg("Invalid number of arguments. Key-value pairs must be even.")
		return
	}

	event := a.logger.Debug()

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			a.logger.Error().Msgf("Invalid key type at index %d. Keys must be strings.", i)
			return
		}

		value := keyValue[i+1]
		setKeyValuesAny(event, key, value)
	}
	event.Msgf(msg)

}

func (a *appLogger) Info(msg string) {
	a.logger.Info().Msg(msg)
}

func (a *appLogger) Infof(msg string, keyValue ...interface{}) {
	if len(keyValue)%2 != 0 {
		a.logger.Error().Msg("Invalid number of arguments. Key-value pairs must be even.")
		return
	}

	event := a.logger.Info()

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			a.logger.Error().Msgf("Invalid key type at index %d. Keys must be strings.", i)
			return
		}

		value := keyValue[i+1]
		setKeyValuesAny(event, key, value)
	}
	event.Msgf(msg)
}

func (a *appLogger) Infop(msg string, args ...interface{}) {
	event := a.logger.Info()

	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if !value.IsValid() {
			continue
		}
		if value.Kind() == reflect.Ptr { //check for ptr
			value = value.Elem()
		}
		key := value.Type().Name()
		setKeyValuesReflect(event, key, value)
	}

	event.Msg(msg)
}

func (a *appLogger) Warn(msg string) {
	a.logger.Warn().Msg(msg)
}

func (a *appLogger) Warnf(msg string, keyValue ...interface{}) {
	if len(keyValue)%2 != 0 {
		a.logger.Error().Msg("Invalid number of arguments. Key-value pairs must be even.")
		return
	}

	event := a.logger.Warn()

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			a.logger.Error().Msgf("Invalid key type at index %d. Keys must be strings.", i)
			return
		}

		value := keyValue[i+1]
		setKeyValuesAny(event, key, value)
	}
	event.Msgf(msg)
}

func (a *appLogger) Error(err error) {
	a.logger.Error().Msg(err.Error())
}

func (a *appLogger) Errorf(msg string, keyValue ...interface{}) {
	if len(keyValue)%2 != 0 {
		a.logger.Error().Msg("Invalid number of arguments. Key-value pairs must be even.")
		return
	}

	event := a.logger.Error()

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			a.logger.Error().Msgf("Invalid key type at index %d. Keys must be strings.", i)
			return
		}

		value := keyValue[i+1]
		setKeyValuesAny(event, key, value)
	}
	event.Msgf(msg)
}

func (a *appLogger) Errorp(msg string, args ...interface{}) {
	event := a.logger.Error()

	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if !value.IsValid() {
			continue
		}
		if value.Kind() == reflect.Ptr { //check for ptr
			value = value.Elem()
		}
		key := value.Type().Name()
		setKeyValuesReflect(event, key, value)
	}

	event.Msg(msg)
}

func (a *appLogger) Panic(msg string) {
	a.logger.Panic().Msg(msg)
}

func (a *appLogger) Panicf(msg string, keyValue ...interface{}) {
	if len(keyValue)%2 != 0 {
		a.logger.Error().Msg("Invalid number of arguments. Key-value pairs must be even.")
		return
	}

	event := a.logger.Panic()

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			a.logger.Error().Msgf("Invalid key type at index %d. Keys must be strings.", i)
			return
		}

		value := keyValue[i+1]
		setKeyValuesAny(event, key, value)
	}
	event.Msgf(msg)
}

func (a *appLogger) Fatal(msg string) {
	a.logger.Fatal().Msg(msg)
}

func (a *appLogger) Fatalf(msg string, keyValue ...interface{}) {
	if len(keyValue)%2 != 0 {
		a.logger.Error().Msg("Invalid number of arguments. Key-value pairs must be even.")
		return
	}

	event := a.logger.Fatal()

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			a.logger.Error().Msgf("Invalid key type at index %d. Keys must be strings.", i)
			return
		}

		value := keyValue[i+1]
		setKeyValuesAny(event, key, value)
	}
	event.Msgf(msg)
}

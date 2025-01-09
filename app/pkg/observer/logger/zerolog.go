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

// Debugp logs warning with pairs of key value
func (a *appLogger) Debugp(msg string, keyValue ...interface{}) {
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

// Infop logs info with pairs of key value
func (a *appLogger) Infop(msg string, keyValue ...interface{}) {
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

// Infoa parses arguments to key value pairs where key is a name of structure field / primitive type and value as value
func (a *appLogger) Infoa(msg string, args ...interface{}) {
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

// Warnp logs warning with pairs of key value
func (a *appLogger) Warnp(msg string, keyValue ...interface{}) {
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

// Error logs error
func (a *appLogger) Error(err error) {
	a.logger.Error().Msg(err.Error())
}

// Errorf logs error using template and args
func (a *appLogger) Errorf(format string, args ...interface{}) {
	a.logger.Error().Msgf(format, args...)
}

// Errorp logs errors with pairs of key value
func (a *appLogger) Errorp(msg string, keyValue ...interface{}) {
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

// Errora parses arguments to key value pairs where key is a name of structure field / primitive type and value as value
func (a *appLogger) Errora(msg string, args ...interface{}) {
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

// Panicp panics with log and pairs of key value
func (a *appLogger) Panicp(msg string, keyValue ...interface{}) {
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

// Fatalp panics with log and pairs of key value
func (a *appLogger) Fatalp(msg string, keyValue ...interface{}) {
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

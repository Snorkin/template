package logger

import (
	"github.com/rs/zerolog"
	"reflect"
)

type LogBuilder struct {
	event *zerolog.Event
}

var Build LogBuilder

type MsgBuilder LogBuilder

func (b LogBuilder) Err() MsgBuilder {
	return MsgBuilder{event: Log.logger.Error()}
}

func (b LogBuilder) Dbg() MsgBuilder {
	return MsgBuilder{event: Log.logger.Debug()}
}

func (b LogBuilder) Wrn() MsgBuilder {
	return MsgBuilder{event: Log.logger.Warn()}
}

func (b LogBuilder) Info() MsgBuilder {
	return MsgBuilder{event: Log.logger.Info()}
}

// Pairs prepares passed key value pairs and passes it forward
func (b MsgBuilder) Pairs(keyValue ...any) MsgBuilder {
	if len(keyValue)%2 != 0 {
		Log.Errorf("Invalid number of arguments %d. Key-value pairs must be even.", len(keyValue))
		return b
	}

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			Log.Errorf("Invalid key type at index %d. Keys must be strings.", i)
			return b
		}

		value := keyValue[i+1]
		setKeyValuesAny(b.event, key, value)
	}
	return b
}

// Args parses args to key (fieldName or primitive type) and value and sets them to event
func (b MsgBuilder) Args(args ...any) MsgBuilder {
	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if !value.IsValid() {
			continue
		}
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		key := value.Type().Name()
		setKeyValuesReflect(b.event, key, value)
	}
	return b
}

// Msg logs message, creates key values of structs and primitive types
func (b MsgBuilder) Msg(msg string) {
	b.event.Msg(msg)
}

// Send logs event without message
func (b MsgBuilder) Send(args ...any) {
	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if !value.IsValid() {
			continue
		}
		if value.Kind() == reflect.Ptr { //check for ptr
			value = value.Elem()
		}
		key := value.Type().Name()
		setKeyValuesReflect(b.event, key, value)
	}

	b.event.Send()
}

// Err logs error message
func (b MsgBuilder) Err(err error) {
	b.event.Err(err).Send()
}

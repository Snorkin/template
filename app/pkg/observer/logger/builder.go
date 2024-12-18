package logger

import (
	"github.com/rs/zerolog"
	"reflect"
)

type LogBuilder struct {
	event *zerolog.Event
}

var Build LogBuilder

type LvlBuilder LogBuilder

type MsgBuilder LogBuilder

func (b LogBuilder) Err() LvlBuilder {
	return LvlBuilder{event: Log.logger.Error()}
}

func (b LogBuilder) Dbg() LvlBuilder {
	return LvlBuilder{event: Log.logger.Debug()}
}

func (b LogBuilder) Wrn() LvlBuilder {
	return LvlBuilder{event: Log.logger.Warn()}
}

func (b LogBuilder) Info() LvlBuilder {
	return LvlBuilder{event: Log.logger.Info()}
}

// Pairs prepares passed key value pairs and passes it forward
func (b LvlBuilder) Pairs(keyValue ...any) MsgBuilder {
	if len(keyValue)%2 != 0 {
		Log.Errorf("Invalid number of arguments. Key-value pairs must be even.")
		return MsgBuilder(b)
	}

	for i := 0; i < len(keyValue); i += 2 {
		key, ok := keyValue[i].(string)
		if !ok {
			Log.Errorf("Invalid key type at index %d. Keys must be strings.", i)
			return MsgBuilder(b)
		}

		value := keyValue[i+1]
		setKeyValuesAny(b.event, key, value)
	}
	return MsgBuilder(b)
}

// Msg logs message, creates key values of structs and primitive types
func (b MsgBuilder) Msg(msg string, args ...any) {
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

	b.event.Msg(msg)
}

// Msg logs message, creates key values of structs and primitive types
func (b LvlBuilder) Msg(msg string, args ...any) {
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

	b.event.Msg(msg)
}

// Msg logs message, creates key values of structs and primitive types
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

// Msg logs message, creates key values of structs and primitive types
func (b LvlBuilder) Send(args ...any) {
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

// Err logs error message
func (b LvlBuilder) Err(err error) {
	b.event.Err(err).Send()
}

package logger

import (
	"example/pkg/observer"
	"fmt"
	"github.com/rs/zerolog"
	"reflect"
	"strings"
	"time"
)

// setKeyValuesAny adds key value pair to event using interface type checking
func setKeyValuesAny(event *zerolog.Event, key string, value any) {
	switch v := value.(type) {
	case string:
		event = event.Str(key, v)
	case int:
		event = event.Int(key, v)
	case bool:
		event = event.Bool(key, v)
	case float64:
		event = event.Float64(key, v)
	default:
		event = event.Interface(key, v)
	}
}

// setKeyValuesReflect adds key value pair to event using reflect type checking
func setKeyValuesReflect(event *zerolog.Event, key string, val reflect.Value) {
	if val.Kind() == reflect.Interface { //check for value if interface is passed
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int64:
		event = event.Int64(key, val.Int())
	case reflect.Float64:
		event = event.Float64(key, val.Float())
	case reflect.String:
		event = event.Str(key, val.String())
	case reflect.Bool:
		event = event.Bool(key, val.Bool())
	case reflect.Struct:
		//complex types handlers (CanInterface is musthave prevents panics)
		if val.CanInterface() && val.Type() == reflect.TypeOf(time.Time{}) { // time.Time
			v := reflect.ValueOf(val.Interface().(time.Time).UTC().String())
			setKeyValuesReflect(event, key, v)
			break
		}

		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if observer.CheckForIgnore(field.Name, field.Tag.Get(observer.TagName)) {
				continue
			}
			key := key + "." + field.Name
			value := val.Field(i)
			setKeyValuesReflect(event, key, value)
		}
	case reflect.Slice:
		var res []string
		for i := 0; i < val.Len(); i++ {
			if !val.Index(i).CanInterface() {
				break
			}
			res = append(res, fmt.Sprintf("%v", val.Index(i).Interface()))
		}
		event = event.Str(key, strings.Join(res, ","))
	case reflect.Map:
		iter := val.MapRange()
		if key != "" {
			key += "."
		}
		for iter.Next() {
			key := key + iter.Key().String()
			value := iter.Value()
			setKeyValuesReflect(event, key, value)
		}
	case reflect.Ptr:
		if !val.IsNil() {
			v := val.Elem()
			setKeyValuesReflect(event, key, v)
		}
	default:
		if !val.CanInterface() {
			break
		}
		if val.IsValid() && val.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if err, ok := val.Interface().(error); ok {
				event = event.Str(key, err.Error())
				break
			}
		}
		event = event.Interface(key, val.Interface())
	}
}

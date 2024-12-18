package errs

import (
	"example/pkg/observer"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// setMapValue sets new entry of reflect.Value to map
func setMapValue(values map[string]any, key string, val reflect.Value) {
	switch val.Kind() {
	case reflect.Int, reflect.Int64:
		values[key] = val.Int()
	case reflect.Float64:
		values[key] = val.Float()
	case reflect.String:
		values[key] = val.String()
	case reflect.Bool:
		values[key] = val.Bool()
	case reflect.Struct:
		//complex types handlers (CanInterface is musthave panics)
		if val.CanInterface() && val.Type() == reflect.TypeOf(time.Time{}) { // time.Time
			v := reflect.ValueOf(val.Interface().(time.Time).UTC().String())
			setMapValue(values, key, v)
			break
		}

		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if observer.CheckForIgnore(field.Name, field.Tag.Get(observer.TagName)) {
				continue
			}
			key := key + "." + field.Name
			value := val.Field(i)
			setMapValue(values, key, value)
		}
	case reflect.Slice:
		var res []string
		for i := 0; i < val.Len(); i++ {
			if !val.Index(i).CanInterface() {
				break
			}
			res = append(res, fmt.Sprintf("%v", val.Index(i).Interface()))
		}
		values[key] = strings.Join(res, ",")
	case reflect.Ptr:
		if !val.IsNil() {
			v := val.Elem()
			setMapValue(values, key, v)
		}
	default:
		values[key] = "unsupported type"
	}
}

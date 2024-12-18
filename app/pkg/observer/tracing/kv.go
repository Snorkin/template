package trace

import (
	"example/pkg/observer"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"strings"
	"time"
)

func setAttr(span trace.Span, key string, val reflect.Value) {
	switch val.Kind() {
	case reflect.Int, reflect.Int64:
		span.SetAttributes(attribute.Int64(key, val.Int()))
	case reflect.Float64:
		span.SetAttributes(attribute.Float64(key, val.Float()))
	case reflect.String:
		span.SetAttributes(attribute.String(key, val.String()))
	case reflect.Bool:
		span.SetAttributes(attribute.Bool(key, val.Bool()))
	case reflect.Struct:
		//complex types handlers (CanInterface is musthave prevents panics)
		if val.CanInterface() && val.Type() == reflect.TypeOf(time.Time{}) { // time.Time
			v := reflect.ValueOf(val.Interface().(time.Time).UTC().String())
			setAttr(span, key, v)
			break
		}

		//struct fields iteration
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if observer.CheckForIgnore(field.Name, field.Tag.Get(observer.TagName)) {
				continue
			}
			key := key + "." + field.Name
			value := val.Field(i)
			setAttr(span, key, value)
		}
	case reflect.Slice:
		var res []string
		for i := 0; i < val.Len(); i++ {
			if !val.Index(i).CanInterface() {
				break
			}
			res = append(res, fmt.Sprintf("%v", val.Index(i).Interface()))
		}
		span.SetAttributes(attribute.String(key, strings.Join(res, ", ")))
	case reflect.Ptr:
		if !val.IsNil() {
			v := val.Elem()
			setAttr(span, key, v)
		}
	default:
		span.SetAttributes(attribute.String(key, "unsupported type"))
	}
}

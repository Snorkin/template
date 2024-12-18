package trace

import (
	"context"
	"example/pkg/observer"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"strings"
	"time"
)

// Start Maps all argument types including structs, slices and primitives. For sensetive info you can use observer:"ignore" tag to not include field in span
// When primitive types passed in it adds span attribute with type of variable as a key and its value as value
func Start(ctx context.Context, name string, args ...any) (context.Context, Span) {
	ctx, span := otel.Tracer("").Start(ctx, name)

	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if value.Kind() == reflect.Ptr { //check for ptr
			value = value.Elem()
		}
		key := value.Type().Name()
		setAttr(span, key, value)
	}

	return ctx, NewSpan(span)
}

func GetTraceIdFromCtx(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

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

type Span struct {
	s trace.Span
}

func NewSpan(s trace.Span) Span {
	return Span{s: s}
}

func (s *Span) End() {
	s.s.End()
}

func (s *Span) Set(key string, val any) {
	v := reflect.ValueOf(val)
	setAttr(s.s, key, v)
}

func (s *Span) GetTraceId() string {
	return s.s.SpanContext().TraceID().String()
}

func (s *Span) Error(err error) error {
	if err == nil {
		return nil
	}
	s.s.SetStatus(codes.Error, fmt.Sprintf("%+v", err))
	s.s.RecordError(err)
	return err
}

func (s *Span) SetName(name string) {
	s.s.SetName(name)
}

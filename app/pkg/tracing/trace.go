package trace

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"strings"
)

func Start(ctx context.Context, name string, args ...any) (context.Context, Span) {
	ctx, span := otel.Tracer("").Start(ctx, name)

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		key := v.Type().Name()
		value := v
		setAttr(span, key, value)
	}

	return ctx, NewSpan(span)
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
		for i := 0; i < val.NumField(); i++ {
			key := key + "." + val.Type().Field(i).Name
			value := val.Field(i)
			setAttr(span, key, value)
		}
	case reflect.Slice:
		var res []string
		for i := 0; i < val.Len(); i++ {
			res = append(res, fmt.Sprintf("%v", val.Index(i).Interface()))
		}
		span.SetAttributes(attribute.String(key, strings.Join(res, ", ")))
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

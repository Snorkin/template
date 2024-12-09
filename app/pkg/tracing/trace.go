package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

func Start(ctx context.Context, name string, args ...any) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(ctx, name)
	context.WithValue(ctx, "traceCtx", span.SpanContext().TraceID().String())

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		key := v.Type().Name()
		value := v
		setAttr(span, key, value)
	}

	return ctx, span
}

func setAttr(span trace.Span, key string, val reflect.Value) {
	switch val.Kind() {
	case reflect.Int:
		span.SetAttributes(attribute.Int(key, int(val.Int())))
	case reflect.Int64:
		span.SetAttributes(attribute.Int64(key, val.Int()))
	case reflect.Float64:
		span.SetAttributes(attribute.Float64(key, val.Float()))
	case reflect.String:
		span.SetAttributes(attribute.String(key, val.String()))
	case reflect.Bool:
		span.SetAttributes(attribute.Bool(key, val.Bool()))
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			key := val.Type().Field(i).Name
			value := val.Field(i)
			setAttr(span, key, value)
		}

	default:
		span.SetAttributes(attribute.String(key, "complex interface"))
	}
}

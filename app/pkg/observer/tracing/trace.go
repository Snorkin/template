package trace

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
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
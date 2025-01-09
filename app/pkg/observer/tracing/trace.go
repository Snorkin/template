package trace

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"runtime"
	"strings"
)

// StartName Creating span with name. Maps all argument types including structs, slices and primitives. For sensitive info you can use observer:"ignore" tag to not include field in span
// When primitive types passed in it adds span attribute with type of variable as a key and its value as value
func StartName(ctx context.Context, name string, args ...any) (context.Context, Span) {
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

// Start Creating span with name as caller. Maps all argument types including structs, slices and primitives. For sensetive info you can use observer:"ignore" tag to not include field in span
// When primitive types passed in it adds span attribute with type of variable as a key and its value as value
func Start(ctx context.Context, args ...any) (context.Context, Span) {
	var name string
	pc, _, _, ok := runtime.Caller(1) // 1 means the caller of this function
	if !ok {
		name = "no caller"
	}
	n := runtime.FuncForPC(pc).Name()
	s := strings.Split(n, "/")
	name = s[len(s)-1]

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

// Set sets new span attribute by key and value
func (s *Span) Set(key string, val any) {
	v := reflect.ValueOf(val)
	setAttr(s.s, key, v)
}

// GetTraceId returns span traceId
func (s *Span) GetTraceId() string {
	return s.s.SpanContext().TraceID().String()
}

// Error sets error to span
func (s *Span) Error(err error) error {
	if err == nil {
		return nil
	}
	s.s.SetStatus(codes.Error, fmt.Sprintf("%+v", err))
	s.s.RecordError(err)
	return err
}

// SetName sets name for span
func (s *Span) SetName(name string) {
	s.s.SetName(name)
}

// Args Maps all argument types including structs, slices and primitives and adds it to existing span. For sensetive info you can use observer:"ignore" tag to not include field in span
// When primitive types passed in it adds span attribute with type of variable as a key and its value as value
func (s *Span) Args(args ...any) {
	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if value.Kind() == reflect.Ptr { //check for ptr
			value = value.Elem()
		}
		key := value.Type().Name()
		setAttr(s.s, key, value)
	}
}

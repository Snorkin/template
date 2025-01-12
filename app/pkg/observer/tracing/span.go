package trace

import (
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

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

// Args Maps all argument types including structs, slices and primitives and adds it to existing span. For sensitive info you can use observer:"ignore" tag to not include field in span
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

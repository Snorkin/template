package trace

import (
	"context"
	"go.opentelemetry.io/otel"
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

// Start Creating span with name as caller. Maps all argument types including structs, slices and primitives. For sensitive info you can use observer:"ignore" tag to not include field in span
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
		if value.Kind() == reflect.Ptr { // get value from pointer
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

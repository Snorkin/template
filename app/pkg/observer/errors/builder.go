package errs

import (
	"errors"
	"example/pkg/observer/logger"
	trace "example/pkg/observer/tracing"
	"fmt"
	"reflect"
	"time"
)

type ErrBuilder struct {
	Errs
	withLog  bool
	withSpan bool
	span     *trace.Span
}

// New inits ErrBuilder with default values
func New() ErrBuilder {
	return ErrBuilder{
		Errs{
			err:        nil,
			msg:        "",
			code:       Internal,
			time:       time.Now().UTC(),
			domain:     "",
			stacktrace: nil,
			values:     make(map[string]any),
		},
		false,
		false,
		nil,
	}
}

// toErrs returns pointer to Errs that implements error interface
func (b ErrBuilder) toErrs() *Errs {
	return &Errs{
		err:        b.err,
		code:       b.code,
		msg:        b.msg,
		domain:     b.domain,
		time:       b.time,
		stacktrace: b.stacktrace,
		traceId:    b.traceId,
		values:     b.values,
	}
}

// Wrap wraps error adding stacktrace and other nice things.
func (b ErrBuilder) Wrap(err error) error {
	if err == nil {
		if b.msg == "" {
			err = errors.New("default error")
		} else {
			err = errors.New(b.msg)
		}
	}

	b.err = err
	b.stacktrace = newStacktrace()

	if b.withSpan {
		b.traceId = b.span.GetTraceId()
		_ = b.span.Error(b.err)
		if b.msg != "" {
			b.span.Args(b.Errs.ToMap())
		}
	}

	if b.withLog {
		logger.Build.Err().Args(b.ToMap()).Err(err)
	}

	return b.toErrs()
}

// ToError returns error interface from ErrBuilder,
// when msg is empty generates default one,
// if err is empty creates new error using msg
func (b ErrBuilder) ToError() error {
	if b.msg == "" {
		b.err = errors.New("default error")
	} else {
		b.err = errors.New(b.msg)
	}
	b.stacktrace = newStacktrace()
	return b.toErrs()
}

// Code sets code on error.
func (b ErrBuilder) Code(code Code) ErrBuilder {
	b.code = code
	return b
}

// In sets domain on error
func (b ErrBuilder) In(domain string) ErrBuilder {
	b.domain = domain
	return b
}

// Msg sets errors msg
func (b ErrBuilder) Msg(msg string, args ...any) ErrBuilder {
	msg = fmt.Sprintf(msg, args...)
	if b.err == nil {
		b.err = errors.New(msg)
	}
	b.msg = msg
	return b
}

// Values parses arguments to key(fieldName/primitiveType) value and saves it in error
func (b ErrBuilder) Values(args ...any) ErrBuilder {
	values := make(map[string]any)
	for _, arg := range args {
		value := reflect.ValueOf(arg)
		if !value.IsValid() {
			continue
		}
		if value.Kind() == reflect.Ptr { //check for ptr
			value = value.Elem()
		}
		key := value.Type().Name()
		setMapValue(values, key, value)
	}
	b.values = values
	return b
}

// Log logs error using info as key and toMap() method as value
func (b ErrBuilder) Log() ErrBuilder {
	b.withLog = true
	return b
}

// Span adds span to builder
func (b ErrBuilder) Span(span trace.Span) ErrBuilder {
	b.withSpan = true
	b.span = &span
	return b
}

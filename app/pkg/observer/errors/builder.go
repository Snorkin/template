package errs

import (
	"errors"
	trace "example/pkg/observer/tracing"
	"fmt"
	"reflect"
	"time"
)

type ErrBuilder Errs

func New() ErrBuilder {
	return ErrBuilder{
		err:        nil,
		msg:        "",
		code:       Internal,
		time:       time.Now().UTC(),
		domain:     "",
		stacktrace: nil,
		values:     make(map[string]any),
	}
}

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
	return b.toErrs()
}

// Wrap wraps error adding stacktrace and other nice things.
func (b ErrBuilder) WrapSpan(err error, span trace.Span) error {
	if err == nil {
		if b.msg == "" {
			err = errors.New("default error")
		} else {
			err = errors.New(b.msg)
		}
	}
	b.err = span.Error(err)
	b.traceId = span.GetTraceId()
	b.stacktrace = newStacktrace()
	return b.toErrs()
}

func (b ErrBuilder) ToError() error {
	if b.msg == "" {
		b.err = errors.New("default error")
	} else {
		b.err = errors.New(b.msg)
	}
	b.stacktrace = newStacktrace()
	return b.toErrs()
}

// Code sets some arbitrary code on errob.
func (b ErrBuilder) Code(code Code) ErrBuilder {
	b.code = code
	return b
}

// In sets some arbitrary domain on errob.
func (b ErrBuilder) In(domain string) ErrBuilder {
	b.domain = domain
	return b
}

func (b ErrBuilder) Msg(msg string, args ...any) ErrBuilder {
	msg = fmt.Sprintf(msg, args...)
	if b.err == nil {
		b.err = errors.New(msg)
	}
	b.msg = msg
	return b
}

// Code sets some arbitrary code on errob.
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
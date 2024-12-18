package errs

import (
	"fmt"
	"time"
)

type ErrBuilder Errs

func New(err error) ErrBuilder {
	return ErrBuilder{
		err:        err,
		msg:        "",
		code:       Internal,
		time:       time.Now().UTC(),
		domain:     "",
		stacktrace: nil,
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
	}
}

// Wrap wraps error adding stacktrace and other nice things.
func (b ErrBuilder) Wrap() error {
	b.stacktrace = newStacktrace()
	return b.toErrs()
}

// Wrapf wraps error adding stacktrace and other nice things.
func (b ErrBuilder) Wrapf(format string, args ...any) error {
	b.msg = fmt.Sprintf(format, args...)
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

func (b ErrBuilder) Msg(msg string) ErrBuilder {
	b.msg = msg
	return b
}

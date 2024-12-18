package errs

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Errs struct {
	err        error
	code       Code
	msg        string
	domain     string
	time       time.Time
	stacktrace *stacktrace
	values     map[string]any
	traceId    string
}

// Msg returns errors msg
func (e *Errs) Msg() string {
	return e.msg
}

// Code returns errors code
func (e *Errs) Code() Code {
	return e.code
}

// Error returns errors text. Implements error interface
func (e *Errs) Error() string {
	return e.err.Error()
}

// AsErrs returns pointer to Errs and bool flag
func AsErrs(err error) (*Errs, bool) {
	var t *Errs
	ok := errors.As(err, &t)
	return t, ok
}

// Unwrap unwraps Errs returning inner error
func (e *Errs) Unwrap() error { return e.err }

// ToMap parses current error state to map[string]any
func (e *Errs) ToMap() map[string]any {
	payload := map[string]any{}

	if err := e.Error(); err != "" {
		payload["error"] = err
	}

	if msg := e.Msg(); msg != "" {
		payload["msg"] = msg
	}

	if traceId := e.traceId; traceId != "" {
		payload["traceId"] = traceId
	}

	if values := e.values; values != nil && len(values) > 0 {
		payload["values"] = values
	}

	if code := e.Code(); code != 0 {
		payload["code"] = code
	}

	if t := e.time; t != (time.Time{}) {
		payload["time"] = t.UTC().String()
	}

	if domain := e.domain; domain != "" {
		payload["domain"] = domain
	}

	if stacktrace := e.Stacktrace(); stacktrace != "" {
		payload["stacktrace"] = stacktrace
	}

	return payload
}

// Stacktrace returns string of errors stacktrace
func (e *Errs) Stacktrace() string {
	var blocks []string
	topFrame := ""

	recursive(e, func(er *Errs) {
		if e.stacktrace != nil && len(e.stacktrace.frames) > 0 {
			var _ string
			if e.err != nil {
				_ = e.err.Error()
			} else {
				_ = ""
			}

			block := e.stacktrace.String(topFrame)

			blocks = append([]string{block}, blocks...)

			topFrame = e.stacktrace.frames[0].String()
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return strings.Join(blocks, "\nThrown: ")
}

// recursive using callback iterates stacktrace
func recursive(err *Errs, tap func(*Errs)) {
	tap(err)

	if err.err == nil {
		return
	}

	if child, ok := AsErrs(err.err); ok {
		recursive(child, tap)
	}
}

// ToJson parses Errs to json format
func (e *Errs) ToJson() ([]byte, error) {
	return json.Marshal(e.ToMap())
}

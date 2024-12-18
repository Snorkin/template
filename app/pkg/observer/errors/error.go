package errs

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Errs struct {
	err        error
	code       Code
	msg        string
	domain     string
	time       time.Time   `observer:"ignore"`
	stacktrace *stacktrace `observer:"ignore"`
}

func (e *Errs) Msg() string {
	return e.msg
}

func (e *Errs) Code() Code {
	return e.code
}

func (e *Errs) Error() string {
	if e.msg != "" {
		return e.msg
	}
	return e.err.Error()
}

func (e *Errs) AsErrs(err error) bool {
	return errors.As(err, &e.err)
}

func ToErrs(err error) (*Errs, bool) {
	var t *Errs
	ok := errors.As(err, &t)
	return t, ok
}

func (e *Errs) Unwrap() error { return e.err }

func (e *Errs) ToMap() map[string]any {
	payload := map[string]any{}

	if err := e.Error(); err != "" {
		payload["error"] = err
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

func (e *Errs) Stacktrace() string {
	var blocks []string
	topFrame := ""

	recursive(e, func(er *Errs) {
		if e.stacktrace != nil && len(e.stacktrace.frames) > 0 {
			var err string
			if e.err != nil {
				err = e.err.Error()
			} else {
				err = ""
			}

			msg := coalesceOrEmpty(e.msg, err, "Error")
			block := fmt.Sprintf("%s\n%s", msg, e.stacktrace.String(topFrame))

			blocks = append([]string{block}, blocks...)

			topFrame = e.stacktrace.frames[0].String()
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return "Error: " + strings.Join(blocks, "\nThrown: ")
}

func recursive(err *Errs, tap func(*Errs)) {
	tap(err)

	if err.err == nil {
		return
	}

	if child, ok := ToErrs(err.err); ok {
		recursive(child, tap)
	}
}

func coalesceOrEmpty[T comparable](v ...T) T {
	result, _ := coalesce(v...)
	return result
}

func coalesce[T comparable](values ...T) (result T, ok bool) {
	for i := range values {
		if values[i] != result {
			result = values[i]
			ok = true
			return
		}
	}

	return
}

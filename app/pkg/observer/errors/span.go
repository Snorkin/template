package errs

import (
	"example/pkg/observer/tracing"
)

func (b ErrBuilder) Trace(span trace.Span) ErrBuilder {
	_ = span.Error(b.err)
	return b
}

package middleware

import (
	"errors"
	errs "example/pkg/observer/errors"
	trace "example/pkg/observer/tracing"
	"github.com/gofiber/fiber/v2"
)

const (
	traceIdHeader = "X-Trace-Id"
)

func (m MdwManager) Start() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := trace.Start(c.Context())
		defer span.End()

		span.SetName(c.Method() + " " + c.Path())

		c.Locals("traceId", ctx)
		c.SetUserContext(ctx)
		c.Response().Header.Set(traceIdHeader, span.GetTraceId())

		span.Set("traceId", span.GetTraceId())
		span.Set("remoteIP", c.IP())
		span.Set("method", c.Method())
		span.Set("host", c.Hostname())
		span.Set("path", c.OriginalURL())
		span.Set("protocol", c.Protocol())
		span.Set("user-agent", c.Get(fiber.HeaderUserAgent))
		span.Set("content-type", c.Get(fiber.HeaderContentType))
		span.Set("fiber-version", fiber.Version)
		span.Set("status-code", c.Response().StatusCode())

		if err := c.Next(); err != nil {
			e, ok := errs.AsErrs(err)
			if !ok {
				errors.As(errs.New().Msg("unwrapped error").Log().Span(span).Wrap(err), &e)
			}
			return e.ToFiberError(c)
		}

		return nil
	}
}

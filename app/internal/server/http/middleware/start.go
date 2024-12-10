package middleware

import (
	"example/pkg/logger"
	"example/pkg/tracing"
	"github.com/gofiber/fiber/v2"
)

const (
	traceIdHeader = "X-Trace-Id"
)

func (m *MdwManager) Start() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := trace.Start(c.Context(), "MdwManager.Error")
		defer span.End()

		// pass the span through userContext
		span.SetName(c.Method() + " " + c.Path())

		c.Locals("traceId", ctx)
		c.SetUserContext(ctx)
		c.Response().Header.Set(traceIdHeader, span.GetTraceId())

		if err := c.Next(); err != nil {
			logger.Log.Errorf("TraceID: %s :%v", span.GetTraceId(), err)
			return span.Error(err)

		}

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

		return nil
	}
}

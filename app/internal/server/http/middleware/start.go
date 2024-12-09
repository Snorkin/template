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

		span.AddAttr("traceId", span.GetTraceId())
		span.AddAttr("remoteIP", c.IP())
		span.AddAttr("method", c.Method())
		span.AddAttr("host", c.Hostname())
		span.AddAttr("path", c.OriginalURL())
		span.AddAttr("protocol", c.Protocol())
		span.AddAttr("user-agent", c.Get(fiber.HeaderUserAgent))
		span.AddAttr("content-type", c.Get(fiber.HeaderContentType))
		span.AddAttr("fiber-version", fiber.Version)
		span.AddAttr("status-code", c.Response().StatusCode())

		return nil
	}
}

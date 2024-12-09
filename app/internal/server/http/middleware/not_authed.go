package middleware

import (
	"example/pkg/tracing"
	"github.com/gofiber/fiber/v2"
)

func (m *MdwManager) NotAuthedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := trace.Start(c.Context(), "MdwManager.NotAuthedMiddleware")
		defer span.End()

		c.Set("traceCtx", span.SpanContext().TraceID().String())
		c.Locals("traceCtx", ctx)

		return c.Next()
	}
}

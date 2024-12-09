package middleware

import (
	"context"
	e "example/pkg/errors"
	"example/pkg/errors/codes"
	"example/pkg/tracing"
	"github.com/gofiber/fiber/v2"
)

func (m *MdwManager) NotAuthedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("traceId").(context.Context)
		if !ok {
			return e.NewCustomError(codes.InvalidArgument, "No trace provided")
		}

		ctx, span := trace.Start(ctx, "MdwManager.NotAuthedMiddleware")
		defer span.End()

		return c.Next()
	}
}

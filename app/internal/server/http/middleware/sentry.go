package middleware

import (
	trace "example/pkg/tracing"
	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

func (m *MdwManager) Sentry() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, span := trace.Start(c.Context(), "MdwManager.Sentry")
		defer span.End()

		err := c.Next()

		if err != nil {
			hub := sentryfiber.GetHubFromContext(c)
			if hub != nil {
				hub.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("path", c.Path())
					scope.SetExtra("request_body", string(c.Body()))
					scope.SetExtra("method", c.Method())
					scope.SetExtra("error", err.Error())
					hub.CaptureException(err)
				})
			}
			return err
		}
		return nil
	}
}

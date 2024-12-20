package errs

import (
	"example/config"
	"github.com/gofiber/fiber/v2"
)

// ToFiberError returns errors state depending on app environment
func (e *Errs) ToFiberError(c *fiber.Ctx) error {
	cfg := config.GetConfig()
	httpCode := e.code.ToHttpCode()

	if cfg.Environment == "dev" {
		return c.Status(httpCode).JSON(fiber.Map{
			"error": e.ToMap(),
		})
	}

	return c.SendStatus(httpCode)
}

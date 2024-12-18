package middleware

import (
	"example/pkg/observer/tracing"
	"github.com/gofiber/fiber/v2"
)

func (m *MdwManager) Start() fiber.Handler {
	return trace.FiberStartMdw()
}

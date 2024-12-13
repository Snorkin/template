package http

import (
	"example/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

// Handler all driver adapters (http/grpc calls, etc)
type httpHandler interface {
	GetUserByLogin() fiber.Handler
	CreateUser() fiber.Handler
}

func MapUserRoutes(userRoutes fiber.Router, mdw middleware.MdwManager, h httpHandler) {
	userRoutes.Get("/:login", mdw.Start(), h.GetUserByLogin())
	userRoutes.Post("/", mdw.Start(), h.CreateUser())
}

package http

import (
	"example/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func MapUserRoutes(userRoutes fiber.Router, mdw middleware.MdwManager, h Handler) {
	userRoutes.Get("/:login", mdw.Start(), h.GetUserByLogin())
	userRoutes.Post("/", mdw.Start(), h.CreateUser())
}

package http

import (
	"example/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func MapUserRoutes(userRoutes fiber.Router, mdw middleware.MdwManager, h Handler) {
	userRoutes.Get("/:login", mdw.NotAuthedMiddleware(), mdw.Sentry(), h.GetUserByLogin())
	userRoutes.Post("/", mdw.NotAuthedMiddleware(), h.CreateUser())
}

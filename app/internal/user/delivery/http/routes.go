package http

import (
	"example/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func MapUserRoutes(userRoutes fiber.Router, mdw middleware.MdwManager, h Handler) {
	userRoutes.Get("/:login", mdw.Start(), mdw.NotAuthedMiddleware(), h.GetUserByLogin())
	userRoutes.Post("/", mdw.Start(), mdw.NotAuthedMiddleware(), h.CreateUser())
}

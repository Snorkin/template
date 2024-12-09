package http

import (
	"example/internal/server/http/middleware"
	"example/internal/user"
	"github.com/gofiber/fiber/v2"
)

func MapUserRoutes(userRoutes fiber.Router, mdw middleware.MdwManager, h user.HttpHandler) {
	userRoutes.Get("/:login", mdw.NotAuthedMiddleware(), h.GetUserByLogin())
	userRoutes.Post("/", mdw.NotAuthedMiddleware(), h.CreateUser())
}

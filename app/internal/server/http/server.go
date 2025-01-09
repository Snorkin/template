package http

import (
	"example/config"
	"example/pkg/observer/logger"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	db   *sqlx.DB
	http *fiber.App
}

func NewServer(
	db *sqlx.DB,
) *Server {
	return &Server{
		http: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		db: db,
	}
}

func (s *Server) Run() error {
	cfg := config.GetConfig()

	s.http.Use(sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	if err := s.MapHandlers(); err != nil {
		logger.Log.Fatalp("Cannot map handlers: %s", err)
	}

	go func() {
		logger.Log.Infoa("HTTP server started", cfg.Server.Http)
		if err := s.http.Listen(cfg.Server.Http.Host + ":" + cfg.Server.Http.Port); err != nil {
			logger.Log.Fatalp("Error starting HTTP server: %s", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() {
	err := s.http.Shutdown()
	if err != nil {
		logger.Log.Errorp("failed to shutdown HTTP server", "err", err)
	} else {
		logger.Log.Infop("HTTP server resolved")
	}
}

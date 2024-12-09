package http

import (
	"example/config"
	"example/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	cfg  config.Config
	log  logger.Logger
	db   *sqlx.DB
	http *fiber.App
}

func NewServer(
	cfg config.Config,
	log logger.Logger,
	db *sqlx.DB,
) *Server {
	return &Server{
		http: fiber.New(fiber.Config{DisableStartupMessage: true}),
		cfg:  cfg,
		log:  log,
		db:   db,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(); err != nil {
		s.log.Fatalf("Cannot map handlers: %s", err)
	}

	go func() {
		s.log.Infof("HTTP server started on: %s:%s", s.cfg.Server.Http.Host, s.cfg.Server.Http.Port)
		if err := s.http.Listen(s.cfg.Server.Http.Host + ":" + s.cfg.Server.Http.Port); err != nil {
			s.log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() {
	err := s.http.Shutdown()
	if err != nil {
		s.log.Error(err)
	} else {
		s.log.Info("HTTP server resolved")
	}
}

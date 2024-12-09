package http

import (
	"example/internal/server/http/middleware"
	userHandler "example/internal/user/delivery/http"
	userRepository "example/internal/user/infra/repository"
	userUsecase "example/internal/user/usecase"
)

func (s *Server) MapHandlers() error {
	//repositories
	userRepo := userRepository.NewUserRepository(s.log, s.db)

	//usecases
	userU := userUsecase.NewUserUsecase(s.cfg, s.log, userRepo)

	//handlers
	userH := userHandler.NewUserHttpHandlers(s.cfg, s.log, userU)

	//middleware
	mdw := middleware.NewMdwManager(s.cfg, s.log)

	//routes
	userR := s.http.Group("user")
	userHandler.MapUserRoutes(userR, mdw, userH)

	return nil
}

package http

import (
	"example/internal/server/http/middleware"
	userHandler "example/internal/user/delivery/http"
	userRepository "example/internal/user/infra/repository"
	userUsecase "example/internal/user/usecase"
)

func (s *Server) MapHandlers() error {
	//repositories
	userRepo := userRepository.NewUserRepository(s.db)

	//usecases
	userU := userUsecase.NewUserUsecase(userRepo)

	//handlers
	userH := userHandler.NewUserHttpHandlers(userU)

	//middleware
	mdw := middleware.NewMdwManager()

	//routes
	userR := s.http.Group("user")
	userHandler.MapUserRoutes(userR, mdw, userH)

	return nil
}

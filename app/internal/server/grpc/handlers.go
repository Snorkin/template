package grpc

import (
	userHandler "example/internal/user/delivery/grpc"
	userRepository "example/internal/user/infra/repository"
	userUsecase "example/internal/user/usecase"
	pbUser "example/pkg/proto"
)

func (s *Server) MapHandlers() {

	//repositories
	userR := userRepository.NewUserRepository(s.db)

	//usecases
	userU := userUsecase.NewUserUsecase(userR)

	//handlers
	userH := userHandler.NewUserHandlers(userU)

	//register handlers
	pbUser.RegisterUserServiceServer(s.grpc, userH)
}

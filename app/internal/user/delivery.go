package user

import (
	pb "example/pkg/proto"
	"github.com/gofiber/fiber/v2"
)

// Handler all driver adapters (http/grpc calls, etc)
type HttpHandler interface {
	GetUserByLogin() fiber.Handler
	CreateUser() fiber.Handler
}

type GrpcHandler interface {
	pb.UserServiceServer
}

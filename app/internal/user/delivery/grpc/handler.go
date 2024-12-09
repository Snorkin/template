package grpc

import (
	"context"
	"example/config"
	"example/internal/user"
	cnv "example/internal/user/delivery/converter/grpc"
	"example/pkg/logger"
	pb "example/pkg/proto"
	"example/pkg/tracing"
)

type Handler struct {
	cfg     config.Config
	usecase user.Usecase
	log     logger.Logger
	pb.UnimplementedUserServiceServer
}

func NewUserHandlers(cfg config.Config, log logger.Logger, usecase user.Usecase) user.GrpcHandler {
	return &Handler{cfg: cfg, log: log, usecase: usecase}
}

func (h Handler) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.GetUserByLoginResponse, error) {
	ctx, span := trace.Start(ctx, "user.GrpcHandler.GetUserByLogin", req)
	defer span.End()

	res, err := h.usecase.GetUserByLogin(ctx, cnv.GetUserByLoginReqDlvrToUc(req))
	if err != nil {
		return nil, err
	}

	return cnv.GetUserByLoginResUcToDlvr(res), nil
}

func (h Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, span := trace.Start(ctx, "user.GrpcHandler.CreateUser", req)
	defer span.End()

	res, err := h.usecase.CreateUser(ctx, cnv.CreateUserReqDlvrToUc(req))
	if err != nil {
		return nil, err
	}

	return cnv.CreateUserResUcToDlvr(res), nil
}

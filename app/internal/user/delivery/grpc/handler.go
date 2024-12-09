package grpc

import (
	"context"
	cnv "example/internal/user/delivery/converter/grpc"
	userUC "example/internal/user/usecase"
	pb "example/pkg/proto"
	"example/pkg/tracing"
)

type Handler interface {
	pb.UserServiceServer
}

type UserGrpcHandler struct {
	usecase userUC.Usecase
	pb.UnimplementedUserServiceServer
}

func NewUserHandlers(usecase userUC.Usecase) *UserGrpcHandler {
	return &UserGrpcHandler{usecase: usecase}
}

func (h UserGrpcHandler) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.GetUserByLoginResponse, error) {
	ctx, span := trace.Start(ctx, "UserGrpcHandler.GetUserByLogin", req)
	defer span.End()

	res, err := h.usecase.GetUserByLogin(ctx, cnv.GetUserByLoginReqDlvrToUc(req))
	if err != nil {
		return nil, err
	}

	return cnv.GetUserByLoginResUcToDlvr(res), nil
}

func (h UserGrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, span := trace.Start(ctx, "UserGrpcHandler.CreateUser", req)
	defer span.End()

	res, err := h.usecase.CreateUser(ctx, cnv.CreateUserReqDlvrToUc(req))
	if err != nil {
		return nil, err
	}

	return cnv.CreateUserResUcToDlvr(res), nil
}

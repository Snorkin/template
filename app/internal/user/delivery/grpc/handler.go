package grpc

import (
	"context"
	cnv "example/internal/user/delivery/converter/grpc"
	uc "example/internal/user/usecase/model"
	"example/pkg/observer/tracing"
	pb "example/pkg/proto"
)

// Usecase main buiseness logic of application
type usecase interface {
	CreateUser(ctx context.Context, req uc.CreateUserReq) (uc.User, error)
	GetUserByLogin(ctx context.Context, req uc.GetUserByLoginReq) (uc.User, error)
}

type UserGrpcHandler struct {
	usecase usecase
	pb.UnimplementedUserServiceServer
}

func NewUserHandlers(usecase usecase) *UserGrpcHandler {
	return &UserGrpcHandler{usecase: usecase}
}

func (h *UserGrpcHandler) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.GetUserByLoginResponse, error) {
	ctx, span := trace.Start(ctx, req)
	defer span.End()

	res, err := h.usecase.GetUserByLogin(ctx, cnv.GetUserByLoginReqDlvrToUc(req))
	if err != nil {
		return nil, err
	}

	return cnv.GetUserByLoginResUcToDlvr(res), nil
}

func (h *UserGrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, span := trace.Start(ctx, req)
	defer span.End()

	res, err := h.usecase.CreateUser(ctx, cnv.CreateUserReqDlvrToUc(req))
	if err != nil {
		return nil, err
	}

	return cnv.CreateUserResUcToDlvr(res), nil
}

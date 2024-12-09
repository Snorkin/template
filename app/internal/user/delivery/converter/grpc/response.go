package converter

import (
	uc "example/internal/user/usecase/model"
	pb "example/pkg/proto"
)

func GetUserByLoginResUcToDlvr(res uc.User) *pb.GetUserByLoginResponse {
	return &pb.GetUserByLoginResponse{
		Id:    res.Id,
		Name:  res.Name,
		Email: res.Email,
		Login: res.Login,
	}
}

func CreateUserResUcToDlvr(res uc.User) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		Id:    res.Id,
		Login: res.Login,
		Name:  res.Name,
		Email: res.Email,
	}
}

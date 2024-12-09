package converter

import (
	uc "example/internal/user/usecase/model"
	pb "example/pkg/proto"
)

func GetUserByLoginReqDlvrToUc(req *pb.GetUserByLoginRequest) uc.GetUserByLoginReq {
	return uc.GetUserByLoginReq{
		Login: req.Login,
	}
}

func CreateUserReqDlvrToUc(req *pb.CreateUserRequest) uc.CreateUserReq {
	return uc.CreateUserReq{
		Name:  req.Name,
		Email: req.Email,
		Login: req.Login,
	}
}

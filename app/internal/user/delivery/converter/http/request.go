package converter

import (
	dlvr "example/internal/user/delivery/model"
	uc "example/internal/user/usecase/model"
)

func GetUserByLoginReqDlvrToUc(req dlvr.GetUserByLoginReq) uc.GetUserByLoginReq {
	return uc.GetUserByLoginReq{
		Login: req.Login,
	}
}

func CreateUserReqDlvrToUc(req dlvr.CreateUserReq) uc.CreateUserReq {
	return uc.CreateUserReq{
		Name:  req.Name,
		Email: req.Email,
		Login: req.Login,
	}
}

package converter

import (
	dlvr "example/internal/user/delivery/model"
	uc "example/internal/user/usecase/model"
)

func GetUserByLoginResUcToDlvr(res uc.User) dlvr.GetUserByLoginRes {
	return dlvr.GetUserByLoginRes{
		Id:    res.Id,
		Name:  res.Name,
		Email: res.Email,
		Login: res.Login,
	}
}

func CreateUserResUcToDlvr(res uc.User) dlvr.CreateUserRes {
	return dlvr.CreateUserRes{
		Id:    res.Id,
		Name:  res.Name,
		Email: res.Email,
		Login: res.Login,
	}
}

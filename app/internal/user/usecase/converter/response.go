package converter

import (
	dlvr "example/internal/user/delivery/model"
	repo "example/internal/user/infra/repository/model"
	uc "example/internal/user/usecase/model"
)

func GetUserByLoginResRepoToUc(res repo.GetUserByLoginRes) uc.User {
	return uc.User{
		Id:    res.Id,
		Name:  res.Name,
		Email: res.Email,
		Login: res.Login,
	}
}

func GetUserByLoginResUcToDlvr(res uc.User) dlvr.GetUserByLoginRes {
	return dlvr.GetUserByLoginRes{
		Id:    res.Id,
		Name:  res.Name,
		Email: res.Email,
	}
}

func CreateUserResRepoToUc(res repo.CreateUserRes, req uc.CreateUserReq) uc.User {
	return uc.User{
		Id:    res.Id,
		Name:  req.Name,
		Email: req.Email,
		Login: req.Login,
	}
}

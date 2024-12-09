package converter

import (
	repo "example/internal/user/infra/repository/model"
	uc "example/internal/user/usecase/model"
)

func GetUserByLoginReqUcToRepo(req uc.GetUserByLoginReq) repo.GetUserByLoginReq {
	return repo.GetUserByLoginReq{
		Login: req.Login,
	}
}

func CreateUserReqUcToRepo(req uc.CreateUserReq) repo.CreateUserReq {
	return repo.CreateUserReq{
		Name:  req.Name,
		Email: req.Email,
		Login: req.Login,
	}
}

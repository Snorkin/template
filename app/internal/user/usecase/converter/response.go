package converter

import (
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

func CreateUserResRepoToUc(res repo.CreateUserRes, req uc.CreateUserReq) uc.User {
	return uc.User{
		Id:    res.Id,
		Name:  req.Name,
		Email: req.Email,
		Login: req.Login,
	}
}

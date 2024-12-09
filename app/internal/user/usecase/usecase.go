package usecase

import (
	"context"
	repo "example/internal/user/infra/repository"
	cnv "example/internal/user/usecase/converter"
	"example/internal/user/usecase/model"
	"example/pkg/tracing"
)

// Usecase main buiseness logic of application
type Usecase interface {
	CreateUser(ctx context.Context, req model.CreateUserReq) (model.User, error)
	GetUserByLogin(ctx context.Context, req model.GetUserByLoginReq) (model.User, error)
}

type UserUsecase struct {
	repo repo.Repository
}

func NewUserUsecase(
	repo repo.Repository,
) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (u UserUsecase) GetUserByLogin(ctx context.Context, req model.GetUserByLoginReq) (model.User, error) {
	ctx, span := trace.Start(ctx, "user.Usecase.GetUserByLogin", req)
	defer span.End()

	res, err := u.repo.GetUserByLogin(ctx, cnv.GetUserByLoginReqUcToRepo(req))
	if err != nil {
		return model.User{}, err
	}

	return cnv.GetUserByLoginResRepoToUc(res), nil
}

func (u UserUsecase) CreateUser(ctx context.Context, req model.CreateUserReq) (model.User, error) {
	ctx, span := trace.Start(ctx, "user.Usecase.CreateUser", req)
	defer span.End()

	res, err := u.repo.CreateUser(ctx, cnv.CreateUserReqUcToRepo(req))
	if err != nil {
		return model.User{}, err
	}

	return cnv.CreateUserResRepoToUc(res, req), nil
}

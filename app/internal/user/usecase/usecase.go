package usecase

import (
	"context"
	repo "example/internal/user/infra/repository/model"
	cnv "example/internal/user/usecase/converter"
	"example/pkg/observer/tracing"

	"example/internal/user/usecase/model"
)

type repositrory interface {
	CreateUser(ctx context.Context, req repo.CreateUserReq) (repo.CreateUserRes, error)
	GetUserByLogin(ctx context.Context, req repo.GetUserByLoginReq) (repo.GetUserByLoginRes, error)
}

type UserUsecase struct {
	repo repositrory
}

func NewUserUsecase(
	repo repositrory,
) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (u UserUsecase) GetUserByLogin(ctx context.Context, req model.GetUserByLoginReq) (model.User, error) {
	ctx, span := trace.Start(ctx, req)
	defer span.End()

	res, err := u.repo.GetUserByLogin(ctx, cnv.GetUserByLoginReqUcToRepo(req))
	if err != nil {
		return model.User{}, err
	}

	return cnv.GetUserByLoginResRepoToUc(res), nil
}

func (u UserUsecase) CreateUser(ctx context.Context, req model.CreateUserReq) (model.User, error) {
	ctx, span := trace.Start(ctx, req)
	defer span.End()

	res, err := u.repo.CreateUser(ctx, cnv.CreateUserReqUcToRepo(req))
	if err != nil {
		return model.User{}, err
	}

	return cnv.CreateUserResRepoToUc(res, req), nil
}

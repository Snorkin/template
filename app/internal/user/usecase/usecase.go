package usecase

import (
	"context"
	"example/config"
	"example/internal/user"
	cnv "example/internal/user/usecase/converter"
	"example/internal/user/usecase/model"
	"example/pkg/logger"
	"example/pkg/tracing"
)

type userUsecase struct {
	cfg  config.Config
	log  logger.Logger
	repo user.Repository
}

func NewUserUsecase(
	cfg config.Config,
	log logger.Logger,
	repo user.Repository,
) user.Usecase {
	return &userUsecase{
		cfg:  cfg,
		log:  log,
		repo: repo,
	}
}

func (u userUsecase) GetUserByLogin(ctx context.Context, req model.GetUserByLoginReq) (model.User, error) {
	ctx, span := trace.Start(ctx, "user.Usecase.GetUserByLogin", req)
	defer span.End()

	res, err := u.repo.GetUserByLogin(ctx, cnv.GetUserByLoginReqUcToRepo(req))
	if err != nil {
		return model.User{}, err
	}

	return cnv.GetUserByLoginResRepoToUc(res), nil
}

func (u userUsecase) CreateUser(ctx context.Context, req model.CreateUserReq) (model.User, error) {
	ctx, span := trace.Start(ctx, "user.Usecase.CreateUser", req)
	defer span.End()

	res, err := u.repo.CreateUser(ctx, cnv.CreateUserReqUcToRepo(req))
	if err != nil {
		return model.User{}, err
	}

	return cnv.CreateUserResRepoToUc(res, req), nil
}

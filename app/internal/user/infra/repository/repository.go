package repository

import (
	"context"
	"example/internal/user"
	model "example/internal/user/infra/repository/model"
	"example/pkg/logger"
	trace "example/pkg/tracing"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	log logger.Logger
	db  *sqlx.DB
}

func NewUserRepository(log logger.Logger, db *sqlx.DB) user.Repository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, req model.GetUserByLoginReq) (model.GetUserByLoginRes, error) {
	ctx, span := trace.Start(ctx, "user.Repository.GetUserByLogin", req)
	defer span.End()

	var res model.GetUserByLoginRes
	err := u.db.GetContext(ctx, &res, queryGetUserById, req.Login)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, req model.CreateUserReq) (model.CreateUserRes, error) {
	ctx, span := trace.Start(ctx, "user.Repository.CreateUser", req)
	defer span.End()

	var res model.CreateUserRes
	err := u.db.GetContext(ctx, &res, queryCreateUser, req.Name, req.Email, req.Login)
	if err != nil {
		return res, err
	}
	return res, nil
}

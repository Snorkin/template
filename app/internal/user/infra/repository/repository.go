package repository

import (
	"context"
	"example/internal/user/infra/repository/model"
	errs "example/pkg/observer/errors"
	"example/pkg/observer/tracing"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, req model.GetUserByLoginReq) (model.GetUserByLoginRes, error) {
	ctx, span := trace.Start(ctx, "user.Repository.GetUserByLogin", req)
	defer span.End()

	var res model.GetUserByLoginRes
	err := u.db.GetContext(ctx, &res, queryGetUserById, req.Login)
	if err != nil {
		return res, errs.New(err).Msg("failed to get user by login").In("User").Trace(span).Wrap()
	}
	return res, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, req model.CreateUserReq) (model.CreateUserRes, error) {
	ctx, span := trace.Start(ctx, "user.Repository.CreateUser", req)
	defer span.End()

	var res model.CreateUserRes
	err := u.db.GetContext(ctx, &res, queryCreateUser, req.Name, req.Email, req.Login)
	if err != nil {
		return res, span.Error(err)
	}
	return res, nil
}

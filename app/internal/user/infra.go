package user

import (
	"context"
	repo "example/internal/user/infra/repository/model"
)

// Infra all driven adapters interfaces (databases, services, api, etc.)
type Infra interface {
	Repository
}

type Repository interface {
	CreateUser(ctx context.Context, req repo.CreateUserReq) (repo.CreateUserRes, error)
	GetUserByLogin(ctx context.Context, req repo.GetUserByLoginReq) (repo.GetUserByLoginRes, error)
}

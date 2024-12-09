package user

import (
	"context"
	uc "example/internal/user/usecase/model"
)

// Usecase main buiseness logic of application
type Usecase interface {
	CreateUser(ctx context.Context, req uc.CreateUserReq) (uc.User, error)
	GetUserByLogin(ctx context.Context, req uc.GetUserByLoginReq) (uc.User, error)
}

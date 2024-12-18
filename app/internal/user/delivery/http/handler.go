package http

import (
	"context"
	"errors"
	cnv "example/internal/user/delivery/converter/http"
	"example/internal/user/delivery/model"
	uc "example/internal/user/usecase/model"
	"example/pkg/observer/tracing"
	"example/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// Usecase main buiseness logic of application
type usecase interface {
	CreateUser(ctx context.Context, req uc.CreateUserReq) (uc.User, error)
	GetUserByLogin(ctx context.Context, req uc.GetUserByLoginReq) (uc.User, error)
}

type UserHttpHandler struct {
	uc usecase
}

func NewUserHttpHandlers(
	uc usecase,
) *UserHttpHandler {
	return &UserHttpHandler{
		uc: uc,
	}
}

func (h UserHttpHandler) GetUserByLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("traceId").(context.Context)
		if !ok {
			return errors.New("no traceId")
		}

		ctx, span := trace.Start(ctx, "user.HttpHandler.GetUserByLogin")
		defer span.End()

		var req model.GetUserByLoginReq
		if err := validator.ReadRequestParam(c, &req); err != nil {
			return err //default errors placeholder
		}

		res, err := h.uc.GetUserByLogin(ctx, cnv.GetUserByLoginReqDlvrToUc(req))
		if err != nil {
			return err
		}

		return c.JSON(cnv.GetUserByLoginResUcToDlvr(res))
	}
}

func (h UserHttpHandler) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("traceId").(context.Context)
		if !ok {
			return errors.New("no traceId")
		}

		ctx, span := trace.Start(ctx, "user.HttpHandler.CreateUser")
		defer span.End()

		var req model.CreateUserReq
		if err := validator.ReadRequestBody(c, &req); err != nil {
			return err //default errors placeholder
		}

		res, err := h.uc.CreateUser(ctx, cnv.CreateUserReqDlvrToUc(req))
		if err != nil {
			return err
		}

		return c.JSON(cnv.CreateUserResUcToDlvr(res))
	}
}

package http

import (
	"context"
	cnv "example/internal/user/delivery/converter/http"
	"example/internal/user/delivery/model"
	user "example/internal/user/usecase"
	e "example/pkg/errors"
	"example/pkg/errors/codes"
	trace "example/pkg/tracing"
	"example/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// Handler all driver adapters (http/grpc calls, etc)
type Handler interface {
	GetUserByLogin() fiber.Handler
	CreateUser() fiber.Handler
}

type UserHttpHandler struct {
	usecase user.Usecase
}

func NewUserHttpHandlers(
	usecase user.Usecase,
) *UserHttpHandler {
	return &UserHttpHandler{
		usecase: usecase,
	}
}

func (h UserHttpHandler) GetUserByLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("traceCtx").(context.Context)
		if !ok {
			return e.NewCustomError(codes.InvalidArgument, "No trace provided")
		}

		ctx, span := trace.Start(ctx, "user.HttpHandler.GetUserByLogin")
		defer span.End()

		var req model.GetUserByLoginReq
		if err := validator.ReadRequestParam(c, &req); err != nil {
			return err //default errors placeholder
		}

		res, err := h.usecase.GetUserByLogin(c.UserContext(), cnv.GetUserByLoginReqDlvrToUc(req))
		if err != nil {
			return err
		}

		return c.JSON(cnv.GetUserByLoginResUcToDlvr(res))
	}
}

func (h UserHttpHandler) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("traceCtx").(context.Context)
		if !ok {
			return e.NewCustomError(codes.InvalidArgument, "No trace provided")
		}

		ctx, span := trace.Start(ctx, "user.HttpHandler.CreateUser")
		defer span.End()

		var req model.CreateUserReq
		if err := validator.ReadRequestBody(c, &req); err != nil {
			return err //default errors placeholder
		}

		res, err := h.usecase.CreateUser(c.UserContext(), cnv.CreateUserReqDlvrToUc(req))
		if err != nil {
			return err
		}

		return c.JSON(cnv.CreateUserResUcToDlvr(res))
	}
}

package http

import (
	"context"
	"example/config"
	"example/internal/user"
	cnv "example/internal/user/delivery/converter/http"
	"example/internal/user/delivery/model"
	e "example/pkg/errors"
	"example/pkg/errors/codes"
	"example/pkg/logger"
	trace "example/pkg/tracing"
	"example/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg     config.Config
	usecase user.Usecase
	log     logger.Logger
}

func NewUserHttpHandlers(
	cfg config.Config,
	log logger.Logger,
	usecase user.Usecase,
) user.HttpHandler {
	return &Handler{
		cfg:     cfg,
		log:     log,
		usecase: usecase,
	}
}

func (h Handler) GetUserByLogin() fiber.Handler {
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

func (h Handler) CreateUser() fiber.Handler {
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

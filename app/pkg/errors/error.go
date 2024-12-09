package errors

import (
	"errors"
	"example/pkg/errors/codes"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type CustomError struct {
	code    codes.ErrorCode
	message string
}

func NewCustomError(code codes.ErrorCode, message string, args ...interface{}) *CustomError {
	msg := fmt.Sprintf(fmt.Sprint("Error: ", message), args...)
	return &CustomError{
		code:    code,
		message: msg,
	}
}

func (e *CustomError) Error() string {
	return e.message
}

func (e *CustomError) Code() *codes.ErrorCode {
	return &e.code
}

func (e *CustomError) ToFiberError(ctx *fiber.Ctx, env string) error {
	httpCode := e.code.ToHttpCode()
	if env == "dev" {
		return ctx.Status(httpCode).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	return ctx.SendStatus(httpCode)
}

func IsCustomError(err error) bool {
	var e *CustomError
	return errors.As(err, &e)
}

func ToCustomError(err error) *CustomError {
	var e *CustomError
	if !errors.As(err, &e) {
		return nil
	}
	return e
}

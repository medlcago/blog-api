package errors

import (
	"blog-api/pkg/response"
	goerrors "errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var fiberErr *fiber.Error
	if goerrors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	var validatorErr validator.ValidationErrors
	if goerrors.As(err, &validatorErr) {
		code = fiber.StatusBadRequest
	}

	var apiErr *Error
	if goerrors.As(err, &apiErr) {
		code = apiErr.Code
	}

	return ctx.Status(code).JSON(response.Response[struct{}]{
		OK:  false,
		Msg: err.Error(),
	})
}

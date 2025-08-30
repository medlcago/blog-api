package errors

import (
	"blog-api/internal/logger"
	"blog-api/pkg/response"
	goerrors "errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/utils/v2"
)

func NewErrorHandler(l *slog.Logger) fiber.ErrorHandler {
	baseLog := l.With(slog.String("component", "errors.ErrorHandler"))

	return func(ctx fiber.Ctx, err error) error {
		reqLog := baseLog.With(
			slog.String(string(logger.RequestIDKey), requestid.FromContext(ctx)),
		)

		code := fiber.StatusInternalServerError
		msg := utils.StatusMessage(code)

		jsonErr := analyzeJSONError(err)
		if jsonErr != nil {
			code = jsonErr.Code
			msg = jsonErr.Message
		}

		var fiberErr *fiber.Error
		if goerrors.As(err, &fiberErr) {
			code = fiberErr.Code
		}

		var validatorErr validator.ValidationErrors
		if goerrors.As(err, &validatorErr) {
			code = fiber.StatusUnprocessableEntity
			msg = utils.StatusMessage(fiber.StatusUnprocessableEntity)
		}

		var apiErr *Error
		if goerrors.As(err, &apiErr) {
			code = apiErr.Code
		}

		if code >= 500 {
			reqLog.Error("internal error",
				slog.Int("status_code", code),
				slog.Any("err", err),
			)
		} else {
			reqLog.Warn("client error",
				slog.Int("status_code", code),
				slog.Any("err", err),
			)
		}

		if code < 500 && msg == utils.StatusMessage(500) {
			msg = err.Error()
		}

		return ctx.Status(code).JSON(response.Response[struct{}]{
			OK:  false,
			Msg: msg,
		})
	}
}

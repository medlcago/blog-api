package errors

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

type JSONError struct {
	Code    int
	Message string
}

func (e JSONError) Error() string {
	return e.Message
}

func analyzeJSONError(err error) *JSONError {
	errMsg := err.Error()

	switch {
	case errMsg == "EOF" || strings.Contains(errMsg, "unexpected end of JSON input"):
		return &JSONError{
			Code:    fiber.StatusBadRequest,
			Message: "Empty or incomplete JSON data",
		}

	case strings.Contains(errMsg, "invalid character"):
		return &JSONError{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid JSON format",
		}

	case strings.Contains(errMsg, "unknown field"):
		return &JSONError{
			Code:    fiber.StatusBadRequest,
			Message: "Unknown field in JSON",
		}

	case strings.Contains(errMsg, "json: cannot unmarshal"):
		return &JSONError{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid data type in JSON",
		}
	}

	return nil
}

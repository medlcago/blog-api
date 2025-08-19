package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	usernameRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{2,30}[A-Za-z0-9]$`)
)

const (
	UsernameTag = "username"
)

func RegisterUsernameValidation(v *validator.Validate) error {
	return v.RegisterValidation(UsernameTag, func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return usernameRe.MatchString(s)
	})
}

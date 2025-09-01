package validator

import "github.com/go-playground/validator/v10"

func New() (*validator.Validate, error) {
	v := validator.New()

	if err := RegisterUsernameValidation(v); err != nil {
		return nil, err
	}

	return v, nil
}

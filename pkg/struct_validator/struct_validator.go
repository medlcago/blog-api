package struct_validator

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.validator.Struct(out)
}

func New(validator *validator.Validate) *StructValidator {
	return &StructValidator{
		validator: validator,
	}
}

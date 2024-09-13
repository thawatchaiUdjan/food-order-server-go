package middlewares

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

func CreateValidator() *StructValidator {
	return &StructValidator{Validator: validator.New()}
}

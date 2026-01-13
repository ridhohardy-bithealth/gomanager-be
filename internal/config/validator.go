package config

import (
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("authaction", func(fl validator.FieldLevel) bool {
		action := fl.Field().String()
		return action == "login" || action == "create"
	})
	return validate
}

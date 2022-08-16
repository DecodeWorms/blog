package pkg

import "github.com/go-playground/validator/v10"

type InitValidator struct {
	*validator.Validate
}

func NewInitValidator() *validator.Validate {
	return validator.New()
}

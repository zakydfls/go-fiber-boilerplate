package validators

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Initialize() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

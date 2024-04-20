package validator

import (
	"regexp"

	v "github.com/go-playground/validator/v10"
)

func IsAlphaNumericType(fl v.FieldLevel) bool {
	value := fl.Field().String()

	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value)
}

package helper

import (
	"sawitpro/constant"
	"unicode"

	"github.com/go-playground/validator"
)

const (
	anyAlphaCapitalTag = "anyAlphaCapital"
	anyNumericTag      = "anyNumeric"
	anySpecialCharTag  = "anySpecialChar"
)

type validatorHelper struct {
	goValidator *validator.Validate
}

func NewValidatorHelper() validatorHelper {
	goValidator := validator.New()

	goValidator.RegisterValidation(anyAlphaCapitalTag, anyAlphaCapital)
	goValidator.RegisterValidation(anyNumericTag, anyNumeric)
	goValidator.RegisterValidation(anySpecialCharTag, anySpecialChar)

	return validatorHelper{
		goValidator: goValidator,
	}
}

func (vald validatorHelper) ValidateStruct(s interface{}) error {
	return vald.goValidator.Struct(s)
}

func anyAlphaCapital(fl validator.FieldLevel) bool {

	str := fl.Field().String()

	for _, char := range str {
		if unicode.IsUpper(char) {
			return true
		}
	}

	return false
}

func anyNumeric(fl validator.FieldLevel) bool {

	str := fl.Field().String()

	for _, char := range str {
		if unicode.IsDigit(char) {
			return true
		}
	}

	return false
}

func anySpecialChar(fl validator.FieldLevel) bool {

	str := fl.Field().String()

	for _, char := range str {
		_, exists := constant.SpecialCharacter[char]
		if exists {
			return true
		}
	}

	return false
}

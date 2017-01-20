package utils

import (
	"reflect"
	"strings"

	"github.com/ttacon/libphonenumber"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidate() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

func IsPhoneNumber(fl validator.FieldLevel) bool {
	phone, err := libphonenumber.Parse(fl.Field().String(), "")
	if err != nil {
		return false
	}
	if !libphonenumber.IsValidNumber(phone) {
		return false
	}
	return true
}

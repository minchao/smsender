package utils

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func IsTimeRFC3339(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	if err != nil {
		return false
	}
	return true
}

func IsTimeUnixMicro(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Int64:
		s := strconv.FormatInt(field.Int(), 10)
		if len(s) == 16 {
			return true
		}
	case reflect.String:
		if matched, _ := regexp.MatchString(`^\d{16}$`, fl.Field().String()); matched {
			return true
		}
	}
	return false
}

func IsRegexp(fl validator.FieldLevel) bool {
	_, err := regexp.Compile(fl.Field().String())
	if err != nil {
		return false
	}
	return true
}

package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/go-playground/validator.v9"
)

func GetInput(body io.Reader, to interface{}, v *validator.Validate) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, to)
	if err != nil {
		return err
	}
	if v != nil {
		if err = v.Struct(to); err != nil {
			return err
		}
	}
	return nil
}

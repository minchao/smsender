package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type errorMessage struct {
	Error            string      `json:"error"`
	ErrorDescription interface{} `json:"error_description,omitempty"`
}

func getInput(body io.Reader, to interface{}, v *validator.Validate) error {
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

func formErrorMessage(err error) errorMessage {
	var (
		e           string = "bad_request"
		description interface{}
	)
	switch err.(type) {
	case validator.ValidationErrors:
		errors := map[string]interface{}{}
		for _, v := range err.(validator.ValidationErrors) {
			errors[v.Field()] = fmt.Sprintf("Invalid validation on field %s: %s", v.Field(), v.Tag())
		}
		description = errors
	default:
		description = err.Error()
	}
	return errorMessage{Error: e, ErrorDescription: description}
}

func render(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func renderInternalServerError(w http.ResponseWriter, err error) error {
	return render(w,
		http.StatusInternalServerError,
		errorMessage{Error: "internal_server_error", ErrorDescription: err.Error()})
}

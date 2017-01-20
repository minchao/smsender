package model

import "net/http"

type Webhook struct {
	Path   string
	Func   func(http.ResponseWriter, *http.Request)
	Method string
}

package model

import "net/http"

// Webhook represents a webhook endpoint.
type Webhook struct {
	Path   string
	Func   func(http.ResponseWriter, *http.Request)
	Method string
}

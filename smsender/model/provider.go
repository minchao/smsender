package model

const NotFoundProvider = "_not_found_"

type Provider interface {
	Name() string
	Send(message *Message, result *MessageResult)
	Callback(register func(webhook *Webhook), receipts chan<- MessageReceipt)
}

type ProviderError struct {
	Error string `json:"error"`
}

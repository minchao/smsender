package model

type Broker interface {
	Name() string
	Send(message *Message, result *MessageResult)
	Callback(register func(webhook *Webhook), receipts chan<- MessageReceipt)
}

type BrokerError struct {
	Error string `json:"error"`
}

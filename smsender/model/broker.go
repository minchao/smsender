package model

type Broker interface {
	Name() string
	Send(message *Message, result *MessageResult)
}

type BrokerError struct {
	Error string `json:"error"`
}

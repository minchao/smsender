package model

type Broker interface {
	Name() string
	Send(msg *Message, result *Result)
}

type BrokerError struct {
	Error string `json:"error"`
}

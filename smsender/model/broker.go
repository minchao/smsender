package model

type Broker interface {
	Name() string
	Send(message *Message, result *MessageResult)
	Callback(receiptsCh chan<- MessageReceipt)
}

type BrokerError struct {
	Error string `json:"error"`
}

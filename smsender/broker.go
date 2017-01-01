package smsender

type Broker interface {
	Name() string
	Send(msg *Message, result *Result)
}

type BrokerError struct {
	Error string `json:"error"`
}

type DummyBroker struct {
	name string
}

func NewDummyBroker(name string) *DummyBroker {
	return &DummyBroker{
		name: name,
	}
}

func (b DummyBroker) Name() string {
	return b.name
}

func (b DummyBroker) Send(msg *Message, result *Result) {
	result.Status = StatusSent.String()
}

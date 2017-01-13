package smsender

import "github.com/minchao/smsender/smsender/model"

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

func (b DummyBroker) Send(msg *model.Message, result *model.Result) {
	result.Status = model.StatusSent.String()
}

package smsender

import (
	log "github.com/Sirupsen/logrus"
)

type Broker interface {
	Name() string
	Send(msg *Message, result *Result)
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
	log.Infof("broker '%s' send message: %+v", b.Name(), *msg)

	result.Status = StatusSent.String()
}

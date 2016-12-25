package smsender

import (
	log "github.com/Sirupsen/logrus"
)

type Broker interface {
	Name() string
	Send(msg Message)
}

type DummyBroker struct {
	name string
}

func NewDummyBroker(name string) DummyBroker {
	return DummyBroker{
		name: name,
	}
}

func (b DummyBroker) Name() string {
	return b.name
}

func (b DummyBroker) Send(msg Message) {
	log.Infof("broker '%s' send message: %+v", b.Name(), msg)
}

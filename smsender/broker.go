package smsender

import (
	log "github.com/Sirupsen/logrus"
)

type Broker interface {
	Name() string
	Send(msg Message)
	Result(c chan Result, r Result)
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

func (b DummyBroker) Send(msg Message) {
	log.Infof("broker '%s' send message: %+v", b.Name(), msg)

	result := NewResult(msg, b)
	result.Status = StatusSent.String()

	b.Result(msg.Result, *result)
}

func (b DummyBroker) Result(c chan Result, r Result) {
	c <- r
	close(c)
}

package smsender

import (
	"fmt"
	"sync"
)

const DefaultBroker = "_default_"

var senderSingleton Sender

type Sender struct {
	routes  []*Route
	brokers map[string]Broker
	in      chan *Message
	mutex   sync.Mutex
	init    sync.Once
}

func SMSender() *Sender {
	senderSingleton.init.Do(func() {
		senderSingleton.brokers = make(map[string]Broker)
		senderSingleton.in = make(chan *Message, 1000)

		senderSingleton.AddBroker(NewDummyBroker(DefaultBroker))
	})

	return &senderSingleton
}

func (s *Sender) AddBroker(broker Broker) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.brokers[broker.Name()]; exists {
		panic(fmt.Sprintf("broker '%s' added", broker.Name()))
	}
	s.brokers[broker.Name()] = broker
}

func (s *Sender) GetBroker(name string) Broker {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if broker, exists := s.brokers[name]; exists {
		return broker
	}
	return nil
}

func (s *Sender) AddRoute(name, pattern, brokerName string) {
	broker := s.GetBroker(brokerName)
	if broker == nil {
		return
	}

	s.routes = append(s.routes, NewRoute(name, pattern, broker))
}

func (s *Sender) GetRoute(name string) *Route {
	for _, route := range s.routes {
		if route.Name == name {
			return route
		}
	}
	return nil
}

func (s *Sender) GetRoutes() []*Route {
	return s.routes
}

func (s *Sender) Stream(from chan *Message) {
	for {
		select {
		case msg := <-from:
			s.in <- msg
		}
	}
}

func (s *Sender) Run() {
	for {
		select {
		case msg := <-s.in:
			go s.walk(msg)
		}
	}
}

func (s *Sender) walk(msg *Message) {
	for _, r := range s.routes {
		if r.Match(msg.Data.To) {
			msg.Route = r.Name
			r.Broker.Send(*msg)
			return
		}
	}
	s.GetBroker(DefaultBroker).Send(*msg)
}

package smsender

import (
	"fmt"
	"sync"
	"errors"
)

const DefaultBroker = "_default_"

var senderSingleton Sender

type Sender struct {
	routes    []*Route
	brokers   map[string]Broker
	in        chan *Message
	out       chan *Message
	workerNum int
	mutex     sync.Mutex
	init      sync.Once
}

func SMSender(workerNum int) *Sender {
	senderSingleton.init.Do(func() {
		senderSingleton.brokers = make(map[string]Broker)
		senderSingleton.in = make(chan *Message, 1000)
		senderSingleton.out = make(chan *Message, 1000)
		senderSingleton.workerNum = workerNum

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

func (s *Sender) AddRoute(route *Route) {
	s.routes = append([]*Route{route}, s.routes...)
}

func (s *Sender) AddRouteWith(name, pattern, brokerName, from string) error {
	broker := s.GetBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	s.AddRoute(NewRoute(name, pattern, broker).SetFrom(from))
	return nil
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
	for i := 0; i < s.workerNum; i++ {
		w := worker{i, s}
		go func(w worker) {
			for msg := range s.out {
				w.process(msg)
			}
		}(w)
	}

	for {
		select {
		case msg := <-s.in:
			s.out <- msg
		}
	}
}

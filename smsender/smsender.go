package smsender

import (
	"errors"
	"fmt"
	"sync"
)

const DefaultBroker = "_default_"

var senderSingleton Sender

type Sender struct {
	router    Router
	brokers   map[string]Broker
	in        chan *Message
	out       chan *Message
	workerNum int
	rwMutex   sync.RWMutex
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

func (s *Sender) GetBroker(name string) Broker {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	if broker, exists := s.brokers[name]; exists {
		return broker
	}
	return nil
}

func (s *Sender) AddBroker(broker Broker) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	if _, exists := s.brokers[broker.Name()]; exists {
		panic(fmt.Sprintf("broker '%s' already added", broker.Name()))
	}
	s.brokers[broker.Name()] = broker
}

func (s *Sender) GetRoutes() []*Route {
	return s.router.GetAll()
}

func (s *Sender) AddRoute(route *Route) {
	s.router.Add(route)
}

func (s *Sender) AddRouteWith(name, pattern, brokerName, from string) error {
	route := s.router.Get(name)
	if route != nil {
		return errors.New("route already exists")
	}
	broker := s.GetBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	s.router.Add(NewRoute(name, pattern, broker).SetFrom(from))
	return nil
}

func (s *Sender) SetRouteWith(name, pattern, brokerName, from string) error {
	broker := s.GetBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	return s.router.Set(name, pattern, broker, from)
}

func (s *Sender) RemoveRoute(name string) {
	s.router.Remove(name)
}

func (s *Sender) ReorderRoutes(rangeStart, rangeLength, insertBefore int) error {
	return s.router.Reorder(rangeStart, rangeLength, insertBefore)
}

func (s *Sender) Match(phone string) (*Route, bool) {
	return s.router.Match(phone)
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
			for message := range s.out {
				w.process(message)
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

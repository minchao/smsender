package smsender

import (
	"errors"
	"fmt"
	"sync"
)

const DefaultBroker = "_default_"

var senderSingleton Sender

type Sender struct {
	Router
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

func (s *Sender) getBroker(name string) Broker {
	if broker, exists := s.brokers[name]; exists {
		return broker
	}
	return nil
}

func (s *Sender) GetBroker(name string) Broker {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.getBroker(name)
}

func (s *Sender) AddRouteWith(name, pattern, brokerName, from string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, route := s.getRoute(name)
	if route != nil {
		return errors.New("route already exists")
	}
	broker := s.getBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	s.AddRoute(NewRoute(name, pattern, broker).SetFrom(from))
	return nil
}

func (s *Sender) SetRouteWith(name, pattern, brokerName, from string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, route := s.getRoute(name)
	if route == nil {
		return errors.New("route not found")
	}
	broker := s.getBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	route.Pattern = pattern
	route.Broker = broker
	route.From = from
	return nil
}

func (s *Sender) ReorderRoutes(rangeStart, rangeLength, insertBefore int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := s.reorder(rangeStart, rangeLength, insertBefore); err != nil {
		return err
	}
	return nil
}

func (s *Sender) RemoveRoute(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	idx, route := s.getRoute(name)
	if route != nil {
		s.removeRoute(idx)
	}
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

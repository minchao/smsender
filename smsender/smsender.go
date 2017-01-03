package smsender

import (
	"errors"
	"fmt"
	"sync"
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

func (s *Sender) AddRoute(route *Route) {
	s.routes = append([]*Route{route}, s.routes...)
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

func (s *Sender) getRoute(name string) (idx int, route *Route) {
	for idx, route := range s.routes {
		if route.Name == name {
			return idx, route
		}
	}
	return idx, nil
}

func (s *Sender) GetRoute(name string) *Route {
	_, route := s.getRoute(name)
	return route
}

func (s *Sender) GetRoutes() []*Route {
	return s.routes
}

func (s *Sender) ReorderRoutes(rangeStart, rangeLength, insertBefore int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	length := len(s.routes)
	if rangeStart > (length - 1) {
		return errors.New("given positions are out of bounds")
	}
	if (rangeStart + rangeLength) > length {
		return errors.New("route selected to be reordered are out of bounds")
	}
	if insertBefore > length {
		return errors.New("given positions are out of bounds")
	}

	var (
		routes   []*Route
		rangeEnd = rangeStart + rangeLength
		subPre   = s.routes[:rangeStart]
		sub      = s.routes[rangeStart:rangeEnd]
		subEnd   = s.routes[rangeEnd:]
	)

	switch {
	case insertBefore == 0:
		routes = append(routes, sub...)
		routes = append(routes, subPre...)
		routes = append(routes, subEnd...)
	case insertBefore < rangeStart:
		routes = append(routes, subPre[:insertBefore]...)
		routes = append(routes, sub...)
		routes = append(routes, subPre[insertBefore:]...)
		routes = append(routes, subEnd...)
	case insertBefore == rangeStart, insertBefore <= rangeEnd:
		routes = s.routes
	case insertBefore < length:
		subBefore := insertBefore - rangeEnd
		routes = append(routes, subPre...)
		routes = append(routes, subEnd[:subBefore]...)
		routes = append(routes, sub...)
		routes = append(routes, subEnd[subBefore:]...)
	case insertBefore == length:
		routes = append(routes, subPre...)
		routes = append(routes, subEnd...)
		routes = append(routes, sub...)
	}

	s.routes = routes

	return nil
}

func (s *Sender) RemoveRoute(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	idx, route := s.getRoute(name)
	if route != nil {
		s.routes = append(s.routes[:idx], s.routes[idx+1:]...)
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

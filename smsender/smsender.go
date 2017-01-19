package smsender

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender/brokers/dummy"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

const DefaultBroker = "_default_"

var senderSingleton Sender

type Sender struct {
	store     store.Store
	router    Router
	brokers   map[string]model.Broker
	in        chan *model.Message
	out       chan *model.Message
	workerNum int
	rwMutex   sync.RWMutex
	init      sync.Once
}

func SMSender(workerNum int) *Sender {
	senderSingleton.init.Do(func() {
		senderSingleton.store = store.NewSqlStore()
		senderSingleton.router = *NewRouter()
		senderSingleton.brokers = make(map[string]model.Broker)
		senderSingleton.in = make(chan *model.Message, 1000)
		senderSingleton.out = make(chan *model.Message, 1000)
		senderSingleton.workerNum = workerNum
		senderSingleton.AddBroker(dummy.NewBroker(DefaultBroker))
	})
	return &senderSingleton
}

func (s *Sender) GetBroker(name string) model.Broker {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	if broker, exists := s.brokers[name]; exists {
		return broker
	}
	return nil
}

func (s *Sender) AddBroker(broker model.Broker) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	if _, exists := s.brokers[broker.Name()]; exists {
		panic(fmt.Sprintf("broker '%s' already added", broker.Name()))
	}
	s.brokers[broker.Name()] = broker
}

func (s *Sender) GetRoutes() []*model.Route {
	return s.router.GetAll()
}

func (s *Sender) AddRoute(route *model.Route) {
	s.router.Add(route)
	s.SaveRoutesToDB()
}

func (s *Sender) AddRouteWith(name, pattern, brokerName, from string, isActive bool) error {
	route := s.router.Get(name)
	if route != nil {
		return errors.New("route already exists")
	}
	broker := s.GetBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	s.router.Add(model.NewRoute(name, pattern, broker, isActive).SetFrom(from))
	s.SaveRoutesToDB()
	return nil
}

func (s *Sender) SetRouteWith(name, pattern, brokerName, from string, isActive bool) error {
	broker := s.GetBroker(brokerName)
	if broker == nil {
		return errors.New("broker not found")
	}
	if err := s.router.Set(name, pattern, broker, from, isActive); err != nil {
		return err
	}
	s.SaveRoutesToDB()
	return nil
}

func (s *Sender) RemoveRoute(name string) {
	s.router.Remove(name)
	s.SaveRoutesToDB()
}

func (s *Sender) ReorderRoutes(rangeStart, rangeLength, insertBefore int) error {
	if err := s.router.Reorder(rangeStart, rangeLength, insertBefore); err != nil {
		return nil
	}
	s.SaveRoutesToDB()
	return nil
}

// Save routes into database.
func (s *Sender) SaveRoutesToDB() error {
	s.router.Lock()
	defer s.router.Unlock()

	var rchan store.StoreChannel

	routes := s.router.getAll()
	rchan = s.store.Route().SaveAll(routes)

	if result := <-rchan; result.Err != nil {
		log.Errorf("SaveRoutesToDB() error: %v", result.Err)
		return result.Err
	}
	return nil
}

// Load routes from database.
func (s *Sender) LoadRoutesFromDB() error {
	var rchan store.StoreChannel

	rchan = s.store.Route().GetAll()

	result := <-rchan
	if result.Err != nil {
		log.Errorf("LoadRoutesFromDB() error: %v", result.Err)
		return result.Err
	}

	routes := []*model.Route{}
	routeRows := result.Data.([]*model.Route)
	for _, r := range routeRows {
		if broker := s.GetBroker(r.Broker); broker != nil {
			routes = append(routes, model.NewRoute(r.Name, r.Pattern, broker, r.IsActive).SetFrom(r.From))
		}
	}

	s.router.SetAll(routes)

	return nil
}

func (s *Sender) GetMessageResults(ids []string) ([]*model.MessageResult, error) {
	if result := <-s.store.Message().GetByIds(ids); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.MessageResult), nil
	}
}

func (s *Sender) Match(phone string) (*model.Route, bool) {
	return s.router.Match(phone)
}

func (s *Sender) GetIncomingQueue() chan *model.Message {
	return s.in
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

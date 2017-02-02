package smsender

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/providers/dummy"
	"github.com/minchao/smsender/smsender/store"
	"github.com/minchao/smsender/smsender/utils"
	config "github.com/spf13/viper"
	"github.com/urfave/negroni"
)

const DefaultProvider = "_default_"

var senderSingleton Sender

type Sender struct {
	store     store.Store
	router    Router
	providers map[string]model.Provider
	in        chan *model.Message
	out       chan *model.Message
	receipts  chan model.MessageReceipt
	workerNum int
	rwMutex   sync.RWMutex
	init      sync.Once

	MuxRouter *mux.Router
}

func SMSender() *Sender {
	senderSingleton.init.Do(func() {
		senderSingleton.store = store.NewSqlStore()
		senderSingleton.router = *NewRouter()
		senderSingleton.providers = make(map[string]model.Provider)
		senderSingleton.in = make(chan *model.Message, 1000)
		senderSingleton.out = make(chan *model.Message, 1000)
		senderSingleton.receipts = make(chan model.MessageReceipt, 1000)
		senderSingleton.workerNum = config.GetInt("worker.num")
		senderSingleton.AddProvider(dummy.NewProvider(DefaultProvider))
		senderSingleton.MuxRouter = mux.NewRouter().StrictSlash(true)
	})
	return &senderSingleton
}

func (s *Sender) GetProvider(name string) model.Provider {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	if provider, exists := s.providers[name]; exists {
		return provider
	}
	return nil
}

func (s *Sender) AddProvider(provider model.Provider) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	if _, exists := s.providers[provider.Name()]; exists {
		panic(fmt.Sprintf("provider '%s' already added", provider.Name()))
	}
	s.providers[provider.Name()] = provider
}

func (s *Sender) GetRoutes() []*model.Route {
	return s.router.GetAll()
}

func (s *Sender) AddRoute(route *model.Route) {
	s.router.Add(route)
	s.SaveRoutesToDB()
}

func (s *Sender) AddRouteWith(name, pattern, providerName, from string, isActive bool) error {
	route := s.router.Get(name)
	if route != nil {
		return errors.New("route already exists")
	}
	provider := s.GetProvider(providerName)
	if provider == nil {
		return errors.New("provider not found")
	}
	s.router.Add(model.NewRoute(name, pattern, provider, isActive).SetFrom(from))
	s.SaveRoutesToDB()
	return nil
}

func (s *Sender) SetRouteWith(name, pattern, providerName, from string, isActive bool) error {
	provider := s.GetProvider(providerName)
	if provider == nil {
		return errors.New("provider not found")
	}
	if err := s.router.Set(name, pattern, provider, from, isActive); err != nil {
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
		if provider := s.GetProvider(r.Provider); provider != nil {
			routes = append(routes, model.NewRoute(r.Name, r.Pattern, provider, r.IsActive).SetFrom(r.From))
		}
	}

	s.router.SetAll(routes)

	return nil
}

func (s *Sender) GetMessageRecords(ids []string) ([]*model.MessageRecord, error) {
	if result := <-s.store.Message().GetByIds(ids); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.MessageRecord), nil
	}
}

func (s *Sender) GetIncomingQueue() chan *model.Message {
	return s.in
}

func (s *Sender) Match(phone string) (*model.Route, bool) {
	return s.router.Match(phone)
}

func (s *Sender) Run() {
	s.initWebhooks()
	s.runWorkers()
	s.runHTTPServer()

	for message := range s.in {
		s.out <- message
	}
}

func (s *Sender) initWebhooks() {
	for _, provider := range s.providers {
		provider.Callback(
			func(webhook *model.Webhook) {
				s.MuxRouter.HandleFunc(webhook.Path, webhook.Func).Methods(webhook.Method)
			},
			s.receipts)
	}
}

func (s *Sender) runWorkers() {
	for i := 0; i < s.workerNum; i++ {
		w := worker{i, s}
		go func(w worker) {
			for {
				select {
				case message := <-s.out:
					w.process(message)
				case receipt := <-s.receipts:
					w.receipt(receipt)
				}
			}
		}(w)
	}
}

func (s *Sender) runHTTPServer() {
	if !config.GetBool("http.enable") {
		return
	}

	n := negroni.New()
	n.UseFunc(utils.Logger)
	n.UseHandler(s.MuxRouter)

	go func() {
		addr := config.GetString("http.addr")
		if config.GetBool("http.tls") {
			log.Infof("Listening for HTTPS on %s", addr)
			log.Fatal(http.ListenAndServeTLS(addr,
				config.GetString("http.tlsCertFile"),
				config.GetString("http.tlsKeyFile"),
				n))
		} else {
			log.Infof("Listening for HTTP on %s", addr)
			log.Fatal(http.ListenAndServe(addr, n))
		}
	}()
}

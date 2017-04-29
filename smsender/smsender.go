package smsender

import (
	"net/http"
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/providers/not_found"
	"github.com/minchao/smsender/smsender/store"
	"github.com/minchao/smsender/smsender/utils"
	config "github.com/spf13/viper"
	"github.com/urfave/negroni"
)

type Sender struct {
	store      store.Store
	messagesCh chan *model.MessageJob
	receiptsCh chan model.MessageReceipt
	workerNum  int

	Router *Router
	// HTTP server router
	HTTPRouter *mux.Router
	siteURL    *url.URL

	shutdown   bool
	shutdownCh chan struct{}
	mutex      sync.RWMutex
	wg         sync.WaitGroup
}

func NewSender() *Sender {
	siteURL, err := url.Parse(config.GetString("http.siteURL"))
	if err != nil {
		log.Fatalln("siteURL err:", err)
	}

	s := store.NewSqlStore()

	return &Sender{
		store:      s,
		messagesCh: make(chan *model.MessageJob, 1000),
		receiptsCh: make(chan model.MessageReceipt, 1000),
		workerNum:  config.GetInt("worker.num"),
		Router:     NewRouter(s, not_found.NewProvider(model.NotFoundProvider)),
		HTTPRouter: mux.NewRouter().StrictSlash(true),
		siteURL:    siteURL,
		shutdownCh: make(chan struct{}, 1),
	}
}

func (s *Sender) SearchMessages(params map[string]interface{}) ([]*model.Message, error) {
	result := <-s.store.Message().Search(params)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.([]*model.Message), nil
}

func (s *Sender) GetMessagesByIds(ids []string) ([]*model.Message, error) {
	result := <-s.store.Message().GetByIds(ids)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.([]*model.Message), nil
}

func (s *Sender) GetMessagesChannel() chan *model.MessageJob {
	return s.messagesCh
}

func (s *Sender) GetSiteURL() *url.URL {
	return s.siteURL
}

// Run performs all startup actions.
func (s *Sender) Run() {
	s.initWebhooks()
	s.initWorkers()
	go s.runHTTPServer()

	select {}
}

// Shutdown sets shutdown flag and stops all workers.
func (s *Sender) Shutdown() {
	s.mutex.Lock()
	if s.shutdown {
		s.mutex.Unlock()
		return
	}
	s.shutdown = true
	s.mutex.Unlock()

	s.wg.Add(s.workerNum)
	close(s.shutdownCh)
	s.wg.Wait()
}

func (s *Sender) initWebhooks() {
	for _, provider := range s.Router.providers {
		provider.Callback(
			func(webhook *model.Webhook) {
				s.HTTPRouter.HandleFunc(webhook.Path, webhook.Func).Methods(webhook.Method)
			},
			s.receiptsCh)
	}
}

func (s *Sender) initWorkers() {
	for i := 0; i < s.workerNum; i++ {
		w := worker{i, s}
		go func(w worker) {
			for {
				select {
				case message := <-s.messagesCh:
					w.process(message)
				case receipt := <-s.receiptsCh:
					w.receipt(receipt)
				case <-s.shutdownCh:
					s.wg.Done()
					return
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
	n.UseHandler(s.HTTPRouter)

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
}

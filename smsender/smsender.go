package smsender

import (
	"net/http"
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

var senderSingleton Sender

type Sender struct {
	store     store.Store
	in        chan *model.Message
	out       chan *model.Message
	receipts  chan model.MessageReceipt
	workerNum int
	init      sync.Once

	Router *Router
	// HTTP server router
	HTTPRouter *mux.Router
}

func SMSender() *Sender {
	senderSingleton.init.Do(func() {
		senderSingleton.store = store.NewSqlStore()
		senderSingleton.in = make(chan *model.Message, 1000)
		senderSingleton.out = make(chan *model.Message, 1000)
		senderSingleton.receipts = make(chan model.MessageReceipt, 1000)
		senderSingleton.workerNum = config.GetInt("worker.num")
		senderSingleton.Router = NewRouter(senderSingleton.store, not_found.NewProvider(model.NotFoundProvider))
		senderSingleton.HTTPRouter = mux.NewRouter().StrictSlash(true)
	})
	return &senderSingleton
}

func (s *Sender) GetMessagesByIds(ids []string) ([]*model.MessageRecord, error) {
	result := <-s.store.Message().GetByIds(ids)
	if result.Err != nil {
		return nil, result.Err
	}
	return result.Data.([]*model.MessageRecord), nil
}

func (s *Sender) GetIncomingQueue() chan *model.Message {
	return s.in
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
	for _, provider := range s.Router.providers {
		provider.Callback(
			func(webhook *model.Webhook) {
				s.HTTPRouter.HandleFunc(webhook.Path, webhook.Func).Methods(webhook.Method)
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
	n.UseHandler(s.HTTPRouter)

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

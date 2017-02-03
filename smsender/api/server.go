package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/model"
	"github.com/rs/cors"
	config "github.com/spf13/viper"
	"github.com/urfave/negroni"
)

type Server struct {
	sender *smsender.Sender
	out    chan<- *model.Message
}

func NewServer(sender *smsender.Sender) *Server {
	server := Server{
		sender: sender,
		out:    sender.GetIncomingQueue(),
	}
	server.init()
	return &server
}

func (s *Server) init() {
	log.Debug("init API")

	router := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)
	router.HandleFunc("/", s.Hello).Methods("GET")
	router.HandleFunc("/routes", s.Routes).Methods("GET")
	router.HandleFunc("/routes", s.RoutePost).Methods("POST")
	router.HandleFunc("/routes", s.RouteReorder).Methods("PUT")
	router.HandleFunc("/routes/{route}", s.RoutePut).Methods("PUT")
	router.HandleFunc("/routes/{route}", s.RouteDelete).Methods("DELETE")
	router.HandleFunc("/routes/test/{phone}", s.RouteTest).Methods("GET")
	router.HandleFunc("/messages/findByIds", s.MessagesFindByIds).Methods("GET")
	router.HandleFunc("/messages", s.MessagesPost).Methods("POST")

	n := negroni.New()

	if config.GetBool("http.api.cors.enable") {
		n.Use(cors.New(cors.Options{
			AllowedOrigins: config.GetStringSlice("http.api.cors.origins"),
			AllowedHeaders: config.GetStringSlice("http.api.cors.headers"),
			AllowedMethods: config.GetStringSlice("http.api.cors.methods"),
			Debug:          config.GetBool("http.api.cors.debug"),
		}))
	}

	n.UseHandler(router)

	s.sender.HTTPRouter.PathPrefix("/api").Handler(n)
}

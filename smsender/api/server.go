package api

import (
	"net/http"

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
	out    chan *model.Message
}

func NewServer(sender *smsender.Sender) *Server {
	server := Server{
		sender: sender,
		out:    make(chan *model.Message, 1000),
	}
	return &server
}

func (s *Server) Run() {
	go s.sender.Stream(s.out)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", s.Hello).Methods("GET")
	r.HandleFunc("/routes", s.Routes).Methods("GET")
	r.HandleFunc("/routes", s.RoutePost).Methods("POST")
	r.HandleFunc("/routes", s.RouteReorder).Methods("PUT")
	r.HandleFunc("/routes/{route}", s.RoutePut).Methods("PUT")
	r.HandleFunc("/routes/{route}", s.RouteDelete).Methods("DELETE")
	r.HandleFunc("/routes/test/{phone}", s.RouteTest).Methods("GET")
	r.HandleFunc("/messages", s.Messages).Methods("GET")
	r.HandleFunc("/messages", s.MessagesPost).Methods("POST")

	n := negroni.New()
	n.UseFunc(logger)

	if config.GetBool("api.cors.enable") {
		n.Use(cors.New(cors.Options{
			AllowedOrigins: config.GetStringSlice("api.cors.origins"),
			AllowedHeaders: config.GetStringSlice("api.cors.headers"),
			AllowedMethods: config.GetStringSlice("api.cors.methods"),
			Debug:          config.GetBool("api.cors.debug"),
		}))
	}

	n.UseHandler(r)

	addr := config.GetString("api.addr")
	if config.GetBool("api.tls") {
		log.Infof("Listening for HTTPS on %s", addr)
		log.Fatal(http.ListenAndServeTLS(addr,
			config.GetString("api.tlsCertFile"),
			config.GetString("api.tlsKeyFile"),
			n))
	} else {
		log.Infof("Listening for HTTP on %s", addr)
		log.Fatal(http.ListenAndServe(addr, n))
	}
}

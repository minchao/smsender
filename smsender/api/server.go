package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender"
	config "github.com/spf13/viper"
	"github.com/urfave/negroni"
)

type Server struct {
	sender *smsender.Sender
	out    chan *smsender.Message
}

func NewServer(sender *smsender.Sender) *Server {
	server := Server{
		sender: sender,
		out:    make(chan *smsender.Message, 1000),
	}
	return &server
}

func (s *Server) Run() {
	go s.sender.Stream(s.out)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", s.Hello).Methods("GET")
	r.HandleFunc("/routes", s.Routes).Methods("GET")
	r.HandleFunc("/send", s.Send).Methods("POST")

	n := negroni.New()
	n.UseFunc(logger)
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

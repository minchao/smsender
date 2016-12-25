package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender"
	"github.com/urfave/negroni"
)

type Message struct {
	Recipient  string `json:"recipient" validate:"required"` // Validate E.164 format
	Body       string `json:"body" validate:"required"`
	Originator string `json:"originator"`
}

type Result struct {
	Message Message `json:"message"`
}

type Server struct {
	addr   string
	sender *smsender.Sender
	in     chan *smsender.Message
}

func NewServer(addr string, sender *smsender.Sender) *Server {
	server := Server{
		addr:   addr,
		sender: sender,
		in:     make(chan *smsender.Message, 1000),
	}
	return &server
}

func (s *Server) Run() {
	go s.sender.Stream(s.in)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", s.Hello).Methods("GET")
	r.HandleFunc("/routes", s.Routes).Methods("GET")
	r.HandleFunc("/send", s.Send).Methods("POST")

	n := negroni.New()
	n.UseFunc(logger)
	n.UseHandler(r)

	log.Infof("Listening for HTTP on %s", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, n))
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	render(w, 200, "Hello!")
}

func (s *Server) Routes(w http.ResponseWriter, r *http.Request) {
	render(w, 200, s.sender.GetRoutes())
}

func (s *Server) Send(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := getInput(r.Body, &msg, newValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	s.in <- &smsender.Message{
		Recipient:  msg.Recipient,
		Body:       msg.Body,
		Originator: msg.Originator,
	}

	// TODO result
	render(w, 200, Result{msg})
}

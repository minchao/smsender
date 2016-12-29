package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender"
	"github.com/rs/xid"
	"github.com/urfave/negroni"
)

type Message struct {
	To   string `json:"to" validate:"required,phone"`
	From string `json:"from"`
	Body string `json:"body" validate:"required"`
}

type Server struct {
	addr   string
	sender *smsender.Sender
	out    chan *smsender.Message
}

func NewServer(addr string, sender *smsender.Sender) *Server {
	server := Server{
		addr:   addr,
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
	var validate = newValidate()
	validate.RegisterValidation("phone", isPhoneNumber)
	err := getInput(r.Body, &msg, validate)
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	resultChan := make(chan smsender.Result, 1)
	s.out <- &smsender.Message{
		Data: smsender.Data{
			Id:   xid.New().String(),
			To:   msg.To,
			From: msg.From,
			Body: msg.Body,
		},
		Result: resultChan,
	}

	result := <-resultChan

	render(w, 200, result)
}

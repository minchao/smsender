package api

import (
	"net/http"

	"github.com/minchao/smsender/smsender"
	"github.com/rs/xid"
)

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	render(w, 200, "Hello!")
}

func (s *Server) Routes(w http.ResponseWriter, r *http.Request) {
	render(w, 200, s.sender.GetRoutes())
}

type Message struct {
	To   string `json:"to" validate:"required,phone"`
	From string `json:"from"`
	Body string `json:"body" validate:"required"`
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

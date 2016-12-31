package api

import (
	"net/http"

	"github.com/minchao/smsender/smsender"
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

	message := smsender.NewMessage(msg.To, msg.From, msg.Body)

	s.out <- message

	result := <-message.Result

	render(w, 200, result)
}

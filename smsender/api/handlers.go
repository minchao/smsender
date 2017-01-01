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
	To    []string `json:"to" validate:"required,gt=0,dive,phone"`
	From  string   `json:"from"`
	Body  string   `json:"body" validate:"required"`
	Async bool     `json:"async,omitempty"`
}

func (s *Server) MessagesPost(w http.ResponseWriter, r *http.Request) {
	var msg Message
	var validate = newValidate()
	validate.RegisterValidation("phone", isPhoneNumber)
	err := getInput(r.Body, &msg, validate)
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	var (
		count         = len(msg.To)
		messageClones = make([]smsender.Message, count)
		resultChans   = make([]<-chan smsender.Result, count)
		results       = make([]smsender.Result, count)
	)

	if count > 100 {
		msg.Async = true
	}

	for i := 0; i < count; i++ {
		message := smsender.NewMessage(msg.To[i], msg.From, msg.Body, msg.Async)
		messageClones[i] = *message
		resultChans[i] = message.Result

		s.out <- message
	}

	if msg.Async {
		for i, message := range messageClones {
			results[i] = *smsender.NewAsyncResult(message)
		}
	} else {
		for i, c := range resultChans {
			results[i] = <-c
		}
	}

	render(w, 200, results)
}

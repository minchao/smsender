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
		resultChans   = make([]<-chan smsender.Result, count)
		messageClones = []smsender.Message{}
		results       = []smsender.Result{}
	)

	if count > 100 {
		msg.Async = true
	}

	for i := 0; i < count; i++ {
		message := smsender.NewMessage(msg.To[i], msg.From, msg.Body, msg.Async)
		resultChans[i] = message.Result
		messageClones = append(messageClones, *message)

		s.out <- message
	}

	if msg.Async {
		for _, message := range messageClones {
			results = append(results, *smsender.NewAsyncResult(message))
		}
	} else {
		for _, c := range resultChans {
			select {
			case result := <-c:
				results = append(results, result)
			}
		}
	}

	render(w, 200, results)
}

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender/model"
)

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	render(w, 200, "Hello!")
}

type RouteResults struct {
	Data []*model.Route `json:"data"`
}

func (s *Server) Routes(w http.ResponseWriter, r *http.Request) {
	render(w, 200, RouteResults{Data: s.sender.GetRoutes()})
}

type Route struct {
	Name     string `json:"name" validate:"required"`
	Pattern  string `json:"pattern" validate:"required"`
	Broker   string `json:"broker" validate:"required"`
	From     string `json:"from"`
	IsActive bool   `json:"is_active"`
}

func (s *Server) RoutePost(w http.ResponseWriter, r *http.Request) {
	var route Route
	err := getInput(r.Body, &route, newValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if err := s.sender.AddRouteWith(route.Name, route.Pattern, route.Broker, route.From, route.IsActive); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, 200, route)
}

type Reorder struct {
	RangeStart   int `json:"range_start" validate:"gte=0"`
	RangeLength  int `json:"range_length" validate:"gte=0"`
	InsertBefore int `json:"insert_before" validate:"gte=0"`
}

func (s *Server) RouteReorder(w http.ResponseWriter, r *http.Request) {
	var reorder Reorder
	err := getInput(r.Body, &reorder, newValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if reorder.RangeLength == 0 {
		reorder.RangeLength = 1
	}
	if err := s.sender.ReorderRoutes(reorder.RangeStart, reorder.RangeLength, reorder.InsertBefore); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, 200, RouteResults{Data: s.sender.GetRoutes()})
}

func (s *Server) RoutePut(w http.ResponseWriter, r *http.Request) {
	var route Route
	err := getInput(r.Body, &route, newValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if err := s.sender.SetRouteWith(route.Name, route.Pattern, route.Broker, route.From, route.IsActive); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, 200, route)
}

func (s *Server) RouteDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeName, _ := vars["route"]
	s.sender.RemoveRoute(routeName)

	render(w, 204, nil)
}

type RouteTestResult struct {
	Phone string       `json:"phone"`
	Route *model.Route `json:"route"`
}

func (s *Server) RouteTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone, _ := vars["phone"]
	validate := newValidate()
	validate.RegisterValidation("phone", isPhoneNumber)
	err := validate.Struct(struct {
		Phone string `json:"phone" validate:"required,phone"`
	}{Phone: phone})
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	route, _ := s.sender.Match(phone)

	render(w, 200, RouteTestResult{Phone: phone, Route: route})
}

type Message struct {
	To    []string `json:"to" validate:"required,gt=0,dive,phone"`
	From  string   `json:"from"`
	Body  string   `json:"body" validate:"required"`
	Async bool     `json:"async,omitempty"`
}

type MessageResults struct {
	Data []model.Result `json:"data"`
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
		messageClones = make([]model.Message, count)
		resultChans   = make([]<-chan model.Result, count)
		results       = make([]model.Result, count)
	)

	if count > 100 {
		msg.Async = true
	}

	for i := 0; i < count; i++ {
		message := model.NewMessage(msg.To[i], msg.From, msg.Body, msg.Async)
		messageClones[i] = *message
		resultChans[i] = message.Result

		s.out <- message
	}

	if msg.Async {
		for i, message := range messageClones {
			results[i] = *model.NewAsyncResult(message)
		}
	} else {
		for i, c := range resultChans {
			results[i] = <-c
		}
	}

	render(w, 200, MessageResults{Data: results})
}

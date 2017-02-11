package api

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/utils"
)

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "Hello!")
}

type provider struct {
	Name string `json:"name"`
}

type routeResults struct {
	Data      []*model.Route `json:"data"`
	Providers []*provider    `json:"providers"`
}

func (s *Server) Routes(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, routeResults{Data: s.sender.Router.GetAll(), Providers: s.getProviders()})
}

type Route struct {
	Name     string `json:"name" validate:"required"`
	Pattern  string `json:"pattern" validate:"required"`
	Provider string `json:"provider" validate:"required"`
	From     string `json:"from"`
	IsActive bool   `json:"is_active"`
}

func (s *Server) RoutePost(w http.ResponseWriter, r *http.Request) {
	var route Route
	err := utils.GetInput(r.Body, &route, utils.NewValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if err := s.sender.Router.AddWith(route.Name, route.Pattern, route.Provider, route.From, route.IsActive); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, http.StatusOK, route)
}

type Reorder struct {
	RangeStart   int `json:"range_start" validate:"gte=0"`
	RangeLength  int `json:"range_length" validate:"gte=0"`
	InsertBefore int `json:"insert_before" validate:"gte=0"`
}

func (s *Server) RouteReorder(w http.ResponseWriter, r *http.Request) {
	var reorder Reorder
	err := utils.GetInput(r.Body, &reorder, utils.NewValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if reorder.RangeLength == 0 {
		reorder.RangeLength = 1
	}
	if err := s.sender.Router.Reorder(reorder.RangeStart, reorder.RangeLength, reorder.InsertBefore); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, http.StatusOK, routeResults{Data: s.sender.Router.GetAll(), Providers: s.getProviders()})
}

func (s *Server) RoutePut(w http.ResponseWriter, r *http.Request) {
	var route Route
	err := utils.GetInput(r.Body, &route, utils.NewValidate())
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if err := s.sender.Router.SetWith(route.Name, route.Pattern, route.Provider, route.From, route.IsActive); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, http.StatusOK, route)
}

func (s *Server) RouteDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeName, _ := vars["route"]
	s.sender.Router.Remove(routeName)

	render(w, http.StatusNoContent, nil)
}

type RouteTestResult struct {
	Phone string       `json:"phone"`
	Route *model.Route `json:"route"`
}

func (s *Server) RouteTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone, _ := vars["phone"]
	validate := utils.NewValidate()
	validate.RegisterValidation("phone", utils.IsPhoneNumber)
	err := validate.Struct(struct {
		Phone string `json:"phone" validate:"required,phone"`
	}{Phone: phone})
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	route, _ := s.sender.Router.Match(phone)

	render(w, http.StatusOK, RouteTestResult{Phone: phone, Route: route})
}

type messagesRequest struct {
	To     string `json:"to" form:"to" validate:"omitempty,phone"`
	Status string `json:"status" form:"status"`
	Since  string `json:"since" form:"since" validate:"omitempty,unixmicro"`
	Until  string `json:"until" form:"until" validate:"omitempty,unixmicro"`
	Limit  int    `json:"limit" form:"limit" validate:"omitempty,gt=0"`
}

type paging struct {
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
}

type ressagesResults struct {
	Data   []*model.MessageRecord `json:"data"`
	Paging paging                 `json:"paging"`
}

func (s *Server) Messages(w http.ResponseWriter, r *http.Request) {
	var req messagesRequest
	if err := form.NewDecoder().Decode(&req, r.URL.Query()); err != nil {
		render(w, http.StatusBadRequest, errorMessage{Error: "bad_request", ErrorDescription: err.Error()})
		return
	}
	if req.Limit == 0 || req.Limit > 100 {
		req.Limit = 100
	}

	validate := utils.NewValidate()
	validate.RegisterValidation("phone", utils.IsPhoneNumber)
	validate.RegisterValidation("unixmicro", utils.IsTimeUnixMicro)
	if err := validate.Struct(req); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	params := make(map[string]interface{})
	if req.To != "" {
		params["to"] = req.To
	}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.Since != "" {
		params["since"], _ = utils.UnixMicroStringToTime(req.Since)
	}
	if req.Until != "" {
		params["until"], _ = utils.UnixMicroStringToTime(req.Until)
	}
	params["limit"] = req.Limit

	messages, err := s.sender.SearchMessages(params)
	if err != nil {
		render(w, http.StatusBadRequest, errorMessage{Error: "not_found", ErrorDescription: err.Error()})
		return
	}

	results := ressagesResults{
		Data:   messages,
		Paging: paging{},
	}

	if len(messages) == 0 {
		results.Data = []*model.MessageRecord{}
	} else {
		// Generate the paging data
		since := messages[0].MessageData.CreatedTime
		until := messages[len(messages)-1].MessageData.CreatedTime

		url, _ := url.Parse("api/messages")
		url = s.sender.GetSiteURL().ResolveReference(url)

		values, _ := form.NewEncoder().Encode(&req)
		cleanEmptyURLValues(&values)

		delete(params, "until")

		params["since"] = since
		prevMessages, _ := s.sender.SearchMessages(params)
		if len(prevMessages) > 0 {
			values.Del("until")
			values.Set("since", strconv.FormatInt(since.UnixNano()/1000, 10))
			url.RawQuery = values.Encode()

			results.Paging.Previous = url.String()
		}

		delete(params, "since")

		params["until"] = until
		nextMessages, _ := s.sender.SearchMessages(params)
		if len(nextMessages) > 0 {
			values.Del("since")
			values.Set("until", strconv.FormatInt(until.UnixNano()/1000, 10))
			url.RawQuery = values.Encode()

			results.Paging.Next = url.String()
		}
	}

	render(w, http.StatusOK, results)
}

type MessagesGetByIdsResults struct {
	Data []*model.MessageRecord `json:"data"`
}

func (s *Server) MessagesGetByIds(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids, _ := r.Form["ids"]
	if err := utils.NewValidate().Struct(struct {
		Ids []string `json:"ids" validate:"required,gt=0,dive,required"`
	}{Ids: ids}); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	messages, err := s.sender.GetMessagesByIds(ids)
	if err != nil {
		render(w, http.StatusNotFound, errorMessage{Error: "not_found", ErrorDescription: err.Error()})
		return
	}
	results := MessagesGetByIdsResults{Data: []*model.MessageRecord{}}
	if len(messages) > 0 {
		results.Data = messages
	}

	render(w, http.StatusOK, results)
}

type MessagesPost struct {
	To    []string `json:"to" validate:"required,gt=0,dive,phone"`
	From  string   `json:"from"`
	Body  string   `json:"body" validate:"required"`
	Async bool     `json:"async,omitempty"`
}

type MessagesPostResults struct {
	Data []model.MessageResult `json:"data"`
}

func (s *Server) MessagesPost(w http.ResponseWriter, r *http.Request) {
	var msg MessagesPost
	var validate = utils.NewValidate()
	validate.RegisterValidation("phone", utils.IsPhoneNumber)
	err := utils.GetInput(r.Body, &msg, validate)
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	var (
		count         = len(msg.To)
		messageClones = make([]model.Message, count)
		resultChans   = make([]<-chan model.MessageResult, count)
		results       = make([]model.MessageResult, count)
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
			results[i] = *model.NewAsyncMessageResult(message)
		}
	} else {
		for i, c := range resultChans {
			results[i] = <-c
		}
	}

	render(w, http.StatusOK, MessagesPostResults{Data: results})
}

func (s *Server) getProviders() []*provider {
	providers := []*provider{}
	for _, p := range s.sender.Router.GetProviders() {
		providers = append(providers, &provider{Name: p.Name()})
	}
	return providers
}

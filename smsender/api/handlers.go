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

type route struct {
	Name     string `json:"name" validate:"required"`
	Pattern  string `json:"pattern" validate:"required,regexp"`
	Provider string `json:"provider" validate:"required"`
	From     string `json:"from"`
	IsActive bool   `json:"is_active"`
}

func (s *Server) RoutePost(w http.ResponseWriter, r *http.Request) {
	var post route
	validate := utils.NewValidate()
	validate.RegisterValidation("regexp", utils.IsRegexp)
	err := utils.GetInput(r.Body, &post, validate)
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if err := s.sender.Router.AddWith(post.Name, post.Pattern, post.Provider, post.From, post.IsActive); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, http.StatusOK, post)
}

type reorder struct {
	RangeStart   int `json:"range_start" validate:"gte=0"`
	RangeLength  int `json:"range_length" validate:"gte=0"`
	InsertBefore int `json:"insert_before" validate:"gte=0"`
}

func (s *Server) RouteReorder(w http.ResponseWriter, r *http.Request) {
	var reorder reorder
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
	var put route
	validate := utils.NewValidate()
	validate.RegisterValidation("regexp", utils.IsRegexp)
	err := utils.GetInput(r.Body, &put, validate)
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}
	if err := s.sender.Router.SetWith(put.Name, put.Pattern, put.Provider, put.From, put.IsActive); err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	render(w, http.StatusOK, put)
}

func (s *Server) RouteDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeName, _ := vars["route"]
	s.sender.Router.Remove(routeName)

	render(w, http.StatusNoContent, nil)
}

type routeTestResult struct {
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

	render(w, http.StatusOK, routeTestResult{Phone: phone, Route: route})
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

type messagesResults struct {
	Data   []*model.Message `json:"data"`
	Paging paging           `json:"paging"`
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

	results := messagesResults{
		Data:   messages,
		Paging: paging{},
	}

	if len(messages) == 0 {
		results.Data = []*model.Message{}
	} else {
		// Generate the paging data
		since := messages[0].CreatedTime
		until := messages[len(messages)-1].CreatedTime

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

type messagess struct {
	Data []*model.Message `json:"data"`
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
	results := messagess{Data: []*model.Message{}}
	if len(messages) > 0 {
		results.Data = messages
	}

	render(w, http.StatusOK, results)
}

type messagesPost struct {
	To    []string `json:"to" validate:"required,gt=0,dive,phone"`
	From  string   `json:"from"`
	Body  string   `json:"body" validate:"required"`
	Async bool     `json:"async,omitempty"`
}

func (s *Server) MessagesPost(w http.ResponseWriter, r *http.Request) {
	var post messagesPost
	var validate = utils.NewValidate()
	validate.RegisterValidation("phone", utils.IsPhoneNumber)
	err := utils.GetInput(r.Body, &post, validate)
	if err != nil {
		render(w, http.StatusBadRequest, formErrorMessage(err))
		return
	}

	var (
		count       = len(post.To)
		jobClones   = make([]*model.MessageJob, count)
		resultChans = make([]<-chan model.Message, count)
		results     = make([]*model.Message, count)
	)

	if count > 100 {
		post.Async = true
	}

	for i := 0; i < count; i++ {
		job := model.NewMessageJob(post.To[i], post.From, post.Body, post.Async)
		jobClones[i] = job
		resultChans[i] = job.Result

		s.out <- job
	}

	if post.Async {
		for i, job := range jobClones {
			results[i] = &job.Message
		}
	} else {
		for i, result := range resultChans {
			message := <-result
			results[i] = &message
		}
	}

	render(w, http.StatusOK, messagess{Data: results})
}

func (s *Server) getProviders() []*provider {
	providers := []*provider{}
	for _, p := range s.sender.Router.GetProviders() {
		providers = append(providers, &provider{Name: p.Name()})
	}
	return providers
}

func (s *Server) Stats(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, model.NewStats())
}

// ShutdownMiddleware will return http.StatusServiceUnavailable if server is already in shutdown progress.
func (s *Server) ShutdownMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if s.sender.IsShutdown() {
		render(w,
			http.StatusServiceUnavailable,
			errorMessage{
				Error:            "service_unavailable",
				ErrorDescription: http.StatusText(http.StatusServiceUnavailable),
			})
		return
	}
	next(w, r)
}

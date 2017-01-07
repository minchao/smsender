package smsender

import (
	"encoding/json"
	"errors"
	"regexp"
)

type Route struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Broker  Broker `json:"broker"`
	From    string `json:"from"`
	regex   *regexp.Regexp
}

func (r *Route) MarshalJSON() ([]byte, error) {
	type Alias Route
	return json.Marshal(&struct {
		*Alias
		Broker string `json:"broker"`
	}{
		Alias:  (*Alias)(r),
		Broker: r.Broker.Name(),
	})
}

func NewRoute(name, pattern string, broker Broker) *Route {
	return &Route{
		Name:    name,
		Pattern: pattern,
		Broker:  broker,
		regex:   regexp.MustCompile(pattern),
	}
}

func (r *Route) SetFrom(from string) *Route {
	r.From = from
	return r
}

func (r *Route) Match(recipient string) bool {
	return r.regex.MatchString(recipient)
}

type Router struct {
	routes []*Route
}

// Add new route to the beginning of routes slice
func (r *Router) AddRoute(route *Route) {
	r.routes = append([]*Route{route}, r.routes...)
}

func (r *Router) getRoute(name string) (idx int, route *Route) {
	for idx, route := range r.routes {
		if route.Name == name {
			return idx, route
		}
	}
	return idx, nil
}

func (r *Router) GetRoute(name string) *Route {
	_, route := r.getRoute(name)
	return route
}

func (r *Router) GetRoutes() []*Route {
	return r.routes
}

func (r *Router) removeRoute(idx int) {
	r.routes = append(r.routes[:idx], r.routes[idx+1:]...)
}

func (r *Router) reorder(rangeStart, rangeLength, insertBefore int) error {
	length := len(r.routes)
	if rangeStart < 0 {
		return errors.New("invalid rangeStart, it should be >= 0")
	}
	if rangeStart > (length - 1) {
		return errors.New("invalid rangeStart, out of bounds")
	}
	if rangeLength <= 0 {
		return errors.New("invalid rangeLength, it should be > 0")
	}
	if (rangeStart + rangeLength) > length {
		return errors.New("route selected to be reordered are out of bounds")
	}
	if insertBefore < 0 {
		return errors.New("invalid insertBefore, it should be >= 0")
	}
	if insertBefore > length {
		return errors.New("invalid insertBefore, out of bounds")
	}

	rangeEnd := rangeStart + rangeLength
	if insertBefore >= rangeStart && insertBefore <= rangeEnd {
		return nil
	}

	var result []*Route

	result = append(result, r.routes[:insertBefore]...)
	result = append(result, r.routes[rangeStart:rangeEnd]...)
	result = append(result, r.routes[insertBefore:]...)
	idxToRemove := rangeStart
	if insertBefore < rangeStart {
		idxToRemove += rangeLength
	}
	result = append(result[:idxToRemove], result[idxToRemove+rangeLength:]...)

	r.routes = result

	return nil
}

func (r *Router) Match(phone string) (*Route, bool) {
	for _, r := range r.routes {
		if r.Match(phone) {
			return r, true
		}
	}
	return nil, false
}

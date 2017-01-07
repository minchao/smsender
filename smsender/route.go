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

func reorderRoutes(routes []*Route, rangeStart, rangeLength, insertBefore int) (result []*Route, err error) {
	length := len(routes)
	if rangeStart < 0 {
		return nil, errors.New("invalid rangeStart, it should be >= 0")
	}
	if rangeStart > (length - 1) {
		return nil, errors.New("invalid rangeStart, out of bounds")
	}
	if rangeLength <= 0 {
		return nil, errors.New("invalid rangeLength, it should be > 0")
	}
	if (rangeStart + rangeLength) > length {
		return nil, errors.New("route selected to be reordered are out of bounds")
	}
	if insertBefore < 0 {
		return nil, errors.New("invalid insertBefore, it should be >= 0")
	}
	if insertBefore > length {
		return nil, errors.New("invalid insertBefore, out of bounds")
	}

	rangeEnd := rangeStart + rangeLength
	if insertBefore >= rangeStart && insertBefore <= rangeEnd {
		result = routes
		return result, nil
	}

	result = append(result, routes[:insertBefore]...)
	result = append(result, routes[rangeStart:rangeEnd]...)
	result = append(result, routes[insertBefore:]...)
	idxToRemove := rangeStart
	if insertBefore < rangeStart {
		idxToRemove += rangeLength
	}
	result = append(result[:idxToRemove], result[idxToRemove+rangeLength:]...)

	return result, nil
}

package smsender

import (
	"encoding/json"
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

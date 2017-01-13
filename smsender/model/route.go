package model

import "regexp"

type Route struct {
	Id      int64  `json:"-"`
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Broker  string `json:"broker"`
	From    string `json:"from" db:"fromName"`
	regex   *regexp.Regexp
	broker  Broker
}

func NewRoute(name, pattern string, broker Broker) *Route {
	return &Route{
		Name:    name,
		Pattern: pattern,
		Broker:  broker.Name(),
		broker:  broker,
		regex:   regexp.MustCompile(pattern),
	}
}

func (r *Route) SetBroker(broker Broker) {
	r.broker = broker
}

func (r *Route) GetBroker() Broker {
	return r.broker
}

func (r *Route) SetFrom(from string) *Route {
	r.From = from
	return r
}

func (r *Route) Match(recipient string) bool {
	return r.regex.MatchString(recipient)
}

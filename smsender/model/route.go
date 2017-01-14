package model

import "regexp"

type Route struct {
	Id       int64  `json:"-"`
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Broker   string `json:"broker"`
	From     string `json:"from" db:"fromName"`
	IsActive bool   `json:"is_active"`
	broker   Broker
	regex    *regexp.Regexp
}

func NewRoute(name, pattern string, broker Broker, isActive bool) *Route {
	return &Route{
		Name:     name,
		Pattern:  pattern,
		Broker:   broker.Name(),
		IsActive: isActive,
		broker:   broker,
		regex:    regexp.MustCompile(pattern),
	}
}

func (r *Route) SetBroker(broker Broker) *Route {
	r.broker = broker
	return r
}

func (r *Route) GetBroker() Broker {
	return r.broker
}

func (r *Route) SetFrom(from string) *Route {
	r.From = from
	return r
}

func (r *Route) Match(recipient string) bool {
	return r.IsActive && r.regex.MatchString(recipient)
}

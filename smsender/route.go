package smsender

import (
	"regexp"
)

type Route struct {
	Name    string
	Pattern string
	Broker  Broker
	regex   *regexp.Regexp
}

func NewRoute(name, pattern string, broker Broker) *Route {
	return &Route{
		Name:    name,
		Pattern: pattern,
		Broker:  broker,
		regex:   regexp.MustCompile(pattern),
	}
}

func (r *Route) Match(recipient string) bool {
	return r.regex.MatchString(recipient)
}

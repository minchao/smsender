package model

import "regexp"

type Route struct {
	Id       int64  `json:"-"`
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Provider string `json:"provider"`
	From     string `json:"from" db:"fromName"`
	IsActive bool   `json:"is_active" db:"isActive"`
	provider Provider
	regex    *regexp.Regexp
}

// NewRoute creates a new instance of the Route.
func NewRoute(name, pattern string, provider Provider, isActive bool) *Route {
	return &Route{
		Name:     name,
		Pattern:  pattern,
		Provider: provider.Name(),
		IsActive: isActive,
		provider: provider,
		regex:    regexp.MustCompile(pattern),
	}
}

func (r *Route) SetPattern(pattern string) *Route {
	r.Pattern = pattern
	r.regex = regexp.MustCompile(pattern)
	return r
}

func (r *Route) SetProvider(provider Provider) *Route {
	r.Provider = provider.Name()
	r.provider = provider
	return r
}

func (r *Route) GetProvider() Provider {
	return r.provider
}

func (r *Route) SetFrom(from string) *Route {
	r.From = from
	return r
}

// Match matches the route against the recipient phone number.
func (r *Route) Match(recipient string) bool {
	return r.IsActive && r.regex.MatchString(recipient)
}

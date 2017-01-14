package smsender

import (
	"errors"
	"sync"

	"github.com/minchao/smsender/smsender/model"
)

// Router registers routes to be matched and dispatches a broker.
type Router struct {
	routes []*model.Route
	sync.RWMutex
}

func (r *Router) getAll() []*model.Route {
	return r.routes
}

func (r *Router) GetAll() []*model.Route {
	r.RLock()
	defer r.RUnlock()
	return r.getAll()
}

func (r *Router) get(name string) (idx int, route *model.Route) {
	for idx, route := range r.routes {
		if route.Name == name {
			return idx, route
		}
	}
	return idx, nil
}

func (r *Router) Get(name string) *model.Route {
	r.RLock()
	defer r.RUnlock()
	_, route := r.get(name)
	return route
}

// Add new route to the beginning of routes slice.
func (r *Router) Add(route *model.Route) {
	r.Lock()
	defer r.Unlock()
	r.routes = append([]*model.Route{route}, r.routes...)
}

func (r *Router) SetAll(routes []*model.Route) {
	r.Lock()
	defer r.Unlock()
	r.routes = routes
}

func (r *Router) Set(name, pattern string, broker model.Broker, from string) error {
	r.Lock()
	defer r.Unlock()
	_, route := r.get(name)
	if route == nil {
		return errors.New("route not found")
	}
	route.Pattern = pattern
	route.Broker = broker.Name()
	route.SetBroker(broker)
	route.From = from
	return nil
}

func (r *Router) Remove(name string) {
	r.Lock()
	defer r.Unlock()
	idx, route := r.get(name)
	if route != nil {
		r.routes = append(r.routes[:idx], r.routes[idx+1:]...)
	}
}

func (r *Router) Reorder(rangeStart, rangeLength, insertBefore int) error {
	r.Lock()
	defer r.Unlock()
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

	var result []*model.Route

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

func (r *Router) Match(phone string) (*model.Route, bool) {
	routes := r.GetAll()
	for _, r := range routes {
		if r.Match(phone) {
			return r, true
		}
	}
	return nil, false
}

package smsender

import (
	"errors"
	"fmt"
	"sync"

	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

// Router registers routes to be matched and dispatches a provider.
type Router struct {
	store     store.Store
	providers map[string]model.Provider
	pMutex    sync.RWMutex
	routes    []*model.Route
	rMutex    sync.RWMutex

	// Configurable Provider to be used when no route matches.
	NotFoundProvider model.Provider
}

// NewRouter creates a new instance of the Router.
func NewRouter(store store.Store, notFoundProvider model.Provider) *Router {
	return &Router{
		store:            store,
		routes:           make([]*model.Route, 0),
		providers:        make(map[string]model.Provider),
		NotFoundProvider: notFoundProvider,
	}
}

func (r *Router) GetProviders() map[string]model.Provider {
	r.pMutex.RLock()
	defer r.pMutex.RUnlock()
	return r.providers
}

func (r *Router) GetProvider(name string) model.Provider {
	r.pMutex.RLock()
	defer r.pMutex.RUnlock()
	if provider, exists := r.providers[name]; exists {
		return provider
	}
	return nil
}

func (r *Router) AddProvider(provider model.Provider) {
	r.pMutex.Lock()
	defer r.pMutex.Unlock()
	if _, exists := r.providers[provider.Name()]; exists {
		panic(fmt.Sprintf("provider '%s' already added", provider.Name()))
	}
	r.providers[provider.Name()] = provider
}

func (r *Router) GetAll() []*model.Route {
	r.rMutex.RLock()
	defer r.rMutex.RUnlock()
	return r.routes
}

func (r *Router) get(name string) (int, *model.Route) {
	for i, route := range r.routes {
		if route.Name == name {
			return i, route
		}
	}
	return 0, nil
}

// Get returns route by specify name.
func (r *Router) Get(name string) *model.Route {
	r.rMutex.RLock()
	defer r.rMutex.RUnlock()
	_, route := r.get(name)
	return route
}

// Add adds new route to the beginning of routes slice.
func (r *Router) Add(route *model.Route) error {
	r.rMutex.Lock()
	defer r.rMutex.Unlock()
	r.routes = append([]*model.Route{route}, r.routes...)
	return r.saveToDB()
}

func (r *Router) AddWith(name, pattern, providerName, from string, isActive bool) error {
	route := r.Get(name)
	if route != nil {
		return errors.New("route already exists")
	}
	provider := r.GetProvider(providerName)
	if provider == nil {
		return errors.New("provider not found")
	}
	return r.Add(model.NewRoute(name, pattern, provider, isActive).SetFrom(from))
}

func (r *Router) Set(name, pattern string, provider model.Provider, from string, isActive bool) error {
	r.rMutex.Lock()
	defer r.rMutex.Unlock()
	_, route := r.get(name)
	if route == nil {
		return errors.New("route not found")
	}
	route.SetPattern(pattern)
	route.SetProvider(provider)
	route.From = from
	route.IsActive = isActive
	return r.saveToDB()
}

func (r *Router) SetWith(name, pattern, providerName, from string, isActive bool) error {
	provider := r.GetProvider(providerName)
	if provider == nil {
		return errors.New("provider not found")
	}
	if err := r.Set(name, pattern, provider, from, isActive); err != nil {
		return err
	}
	return nil
}

func (r *Router) Remove(name string) error {
	r.rMutex.Lock()
	defer r.rMutex.Unlock()
	idx, route := r.get(name)
	if route != nil {
		r.routes = append(r.routes[:idx], r.routes[idx+1:]...)
	}
	return r.saveToDB()
}

func (r *Router) Reorder(rangeStart, rangeLength, insertBefore int) error {
	r.rMutex.Lock()
	defer r.rMutex.Unlock()
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
	return r.saveToDB()
}

func (r *Router) saveToDB() error {
	if result := <-r.store.Route().SaveAll(r.routes); result.Err != nil {
		return result.Err
	}
	return nil
}

// SaveToDB saves the routes into database.
func (r *Router) SaveToDB() error {
	r.rMutex.RLock()
	defer r.rMutex.RUnlock()
	return r.saveToDB()
}

// LoadFromDB returns the routes from database.
func (r *Router) LoadFromDB() error {
	r.rMutex.Lock()
	defer r.rMutex.Unlock()
	result := <-r.store.Route().GetAll()
	if result.Err != nil {
		return result.Err
	}

	routes := []*model.Route{}
	routeRows := result.Data.([]*model.Route)
	for _, row := range routeRows {
		if provider := r.GetProvider(row.Provider); provider != nil {
			routes = append(routes, model.NewRoute(row.Name, row.Pattern, provider, row.IsActive).SetFrom(row.From))
		}
	}

	r.routes = routes

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

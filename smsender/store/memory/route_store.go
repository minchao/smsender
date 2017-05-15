package memory

import (
	"sync"

	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

type RouteStore struct {
	*Store
	routes []*model.Route
	sync.RWMutex
}

func NewMemoryRouteStore(memoryStore *Store) store.RouteStore {
	return &RouteStore{memoryStore, []*model.Route{}, sync.RWMutex{}}
}

func (s *RouteStore) GetAll() store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		result.Data = s.routes

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *RouteStore) SaveAll(routes []*model.Route) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		s.routes = routes

		result.Data = routes

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

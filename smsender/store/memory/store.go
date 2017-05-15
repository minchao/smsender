package memory

import (
	"github.com/minchao/smsender/smsender/plugin"
	"github.com/minchao/smsender/smsender/store"
	"github.com/spf13/viper"
)

func init() {
	plugin.RegisterStore("memory", Plugin)
}

func Plugin(config *viper.Viper) (store.Store, error) {
	return New(), nil
}

type Store struct {
	route   store.RouteStore
	message store.MessageStore
}

func New() store.Store {
	memoryStore := &Store{}
	memoryStore.route = NewMemoryRouteStore(memoryStore)
	memoryStore.message = NewMemoryMessageStore(memoryStore)

	return memoryStore
}

func (s *Store) Route() store.RouteStore {
	return s.route
}

func (s *Store) Message() store.MessageStore {
	return s.message
}

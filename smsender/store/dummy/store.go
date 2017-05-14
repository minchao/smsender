package dummy

import "github.com/minchao/smsender/smsender/store"

type Store struct {
	DummyRoute   store.RouteStore
	DummyMessage store.MessageStore
}

func (s *Store) Route() store.RouteStore {
	return s.DummyRoute
}

func (s *Store) Message() store.MessageStore {
	return s.DummyMessage
}

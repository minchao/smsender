package store

type DummyStore struct {
	DummyRoute   RouteStore
	DummyMessage MessageStore
}

func (ds *DummyStore) Route() RouteStore {
	return ds.DummyRoute
}

func (ds *DummyStore) Message() MessageStore {
	return ds.DummyMessage
}

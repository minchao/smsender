package store

import "github.com/minchao/smsender/smsender/model"

type DummyRouteStore struct {
}

func (rs *DummyRouteStore) GetAll() StoreChannel {
	return make(StoreChannel, 1)
}

func (rs *DummyRouteStore) SaveAll(routes []*model.Route) StoreChannel {
	return make(StoreChannel, 1)
}

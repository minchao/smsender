package dummy

import (
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

type RouteStore struct {
}

func (rs *RouteStore) GetAll() store.Channel {
	return make(store.Channel, 1)
}

func (rs *RouteStore) SaveAll(routes []*model.Route) store.Channel {
	return make(store.Channel, 1)
}

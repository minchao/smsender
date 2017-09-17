package store

import (
	"github.com/minchao/smsender/smsender/model"
)

type Result struct {
	Data interface{}
	Err  error
}

type Channel chan Result

type Store interface {
	Route() RouteStore
	Message() MessageStore
}

type RouteStore interface {
	GetAll() Channel
	SaveAll(routes []*model.Route) Channel
}

type MessageStore interface {
	Get(id string) Channel
	GetByIds(ids []string) Channel
	GetByProviderAndMessageID(provider, providerMessageID string) Channel
	Save(message *model.Message) Channel
	Search(params map[string]interface{}) Channel
	Update(message *model.Message) Channel
}

package store

import "github.com/minchao/smsender/smsender/model"

type StoreResult struct {
	Data interface{}
	Err  error
}

type StoreChannel chan StoreResult

type Store interface {
	Route() RouteStore
	Message() MessageStore
}

type RouteStore interface {
	FindAll() StoreChannel
	SaveAll(routes []*model.Route) StoreChannel
}

type MessageStore interface {
	Find(id string) StoreChannel
	Save(message *model.Result) StoreChannel
	Update(message *model.Result) StoreChannel
}

package store

import "github.com/minchao/smsender/smsender/model"

type StoreResult struct {
	Data interface{}
	Err  error
}

type StoreChannel chan StoreResult

type Store interface {
	Route() RouteStore
}

type RouteStore interface {
	FindAll() StoreChannel
	SaveAll(routes []*model.Route) StoreChannel
}

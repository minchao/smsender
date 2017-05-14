package sql

import (
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

const SqlRouteTable = `
CREATE TABLE IF NOT EXISTS route (
  id       int(11) NOT NULL AUTO_INCREMENT,
  name     varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  pattern  varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  provider varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  fromName varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  isActive tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (id),
  UNIQUE KEY name (name)
) DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`

type RouteStore struct {
	*Store
}

func NewSqlRouteStore(sqlStore *Store) store.RouteStore {
	rs := &RouteStore{sqlStore}

	rs.db.MustExec(SqlRouteTable)

	return rs
}

func (rs *RouteStore) GetAll() store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}
		var routes []*model.Route
		if err := rs.db.Select(&routes, `SELECT * FROM route ORDER BY id ASC`); err != nil {
			result.Err = err
		} else {
			result.Data = routes
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (rs *RouteStore) SaveAll(routes []*model.Route) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		tx := rs.db.MustBegin()
		tx.MustExec(`TRUNCATE TABLE route`)
		for _, route := range routes {
			tx.MustExec(`INSERT INTO route
				(name, pattern, provider, fromName, isActive)
				VALUES (?, ?, ?, ?, ?)`,
				route.Name, route.Pattern, route.Provider, route.From, route.IsActive)
		}
		if err := tx.Commit(); err != nil {
			result.Err = err
		} else {
			result.Data = routes
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

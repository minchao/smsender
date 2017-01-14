package store

import "github.com/minchao/smsender/smsender/model"

const SqlRouteTable = `
CREATE TABLE IF NOT EXISTS route (
  id       int(11) NOT NULL AUTO_INCREMENT,
  name     varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  pattern  varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  broker   varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  fromName varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  isActive tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (id),
  UNIQUE KEY name (name)
) DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`

type SqlRouteStore struct {
	*SqlStore
}

func NewSqlRouteStore(sqlStore *SqlStore) RouteStore {
	rs := &SqlRouteStore{sqlStore}

	rs.db.MustExec(SqlRouteTable)

	return rs
}

func (rs *SqlRouteStore) FindAll() StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}
		var routes []*model.Route
		if err := rs.db.Select(&routes, `SELECT * FROM route`); err != nil {
			result.Err = err
		} else {
			result.Data = routes
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (rs *SqlRouteStore) SaveAll(routes []*model.Route) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		tx := rs.db.MustBegin()
		tx.MustExec(`TRUNCATE TABLE route`)
		for _, route := range routes {
			tx.MustExec(`INSERT INTO route
				(name, pattern, broker, fromName, isActive)
				VALUES (?, ?, ?, ?, ?)`,
				route.Name, route.Pattern, route.Broker, route.From, route.IsActive)
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

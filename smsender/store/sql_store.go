package store

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	config "github.com/spf13/viper"
)

type SqlStore struct {
	db      *sqlx.DB
	route   RouteStore
	message MessageStore
}

func initConnection() *SqlStore {
	sqlStore := &SqlStore{}

	db, err := sqlx.Connect("mysql", config.GetString("db.dsn"))
	if err != nil {
		log.Fatalf("initDB error: %v", err)
	}
	db.SetMaxOpenConns(config.GetInt("db.connection.maxOpenConns"))
	db.SetMaxIdleConns(config.GetInt("db.connection.maxIdleConns"))

	sqlStore.db = db

	return sqlStore
}

func NewSqlStore() Store {
	sqlStore := initConnection()

	sqlStore.route = NewSqlRouteStore(sqlStore)
	sqlStore.message = NewSqlMessageStore(sqlStore)

	return sqlStore
}

func (ss *SqlStore) Route() RouteStore {
	return ss.route
}

func (ss *SqlStore) Message() MessageStore {
	return ss.message
}

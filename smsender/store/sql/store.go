package sql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/minchao/smsender/smsender/plugin"
	"github.com/minchao/smsender/smsender/store"
	"github.com/spf13/viper"
)

func init() {
	plugin.RegisterStore("sql", Plugin)
}

func Plugin(config *viper.Viper) (store.Store, error) {
	return New(config)
}

type Store struct {
	db      *sqlx.DB
	route   store.RouteStore
	message store.MessageStore
}

func New(config *viper.Viper) (store.Store, error) {
	sqlStore := &Store{}

	db, err := sqlx.Connect(config.GetString("driver"), config.GetString("dsn"))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.GetInt("connection.maxOpenConns"))
	db.SetMaxIdleConns(config.GetInt("connection.maxIdleConns"))

	sqlStore.db = db

	sqlStore.route = NewSQLRouteStore(sqlStore)
	sqlStore.message = NewSQLMessageStore(sqlStore)

	return sqlStore, nil
}

func (s *Store) Route() store.RouteStore {
	return s.route
}

func (s *Store) Message() store.MessageStore {
	return s.message
}

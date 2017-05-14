package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

const SqlMessageTable = `
CREATE TABLE IF NOT EXISTS message (
  id                varchar(40) COLLATE utf8_unicode_ci NOT NULL,
  toNumber          varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  fromName          varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  body              text COLLATE utf8_unicode_ci NOT NULL,
  async             tinyint(1) NOT NULL DEFAULT '0',
  route             varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  provider          varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  providerMessageId varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  steps             json DEFAULT NULL,
  status            varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  createdTime       datetime(6) NOT NULL,
  updatedTime       datetime(6) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY providerMessageId (provider, providerMessageId)
) DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`

type MessageStore struct {
	*Store
}

func NewSqlMessageStore(sqlStore *Store) store.MessageStore {
	ms := &MessageStore{sqlStore}

	ms.db.MustExec(SqlMessageTable)

	return ms
}

func (ms *MessageStore) Get(id string) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		var message model.Message
		if err := ms.db.Get(&message, `SELECT * FROM message WHERE id = ?`, id); err != nil {
			result.Err = err
		} else {
			result.Data = &message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *MessageStore) GetByIds(ids []string) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		query, args, err := sqlx.In("SELECT * FROM message WHERE id IN (?)", ids)
		if err != nil {
			result.Err = err
			storeChannel <- result
			close(storeChannel)
			return
		}
		query = ms.db.Rebind(query)

		var messages []*model.Message
		if err := ms.db.Select(&messages, query, args...); err != nil {
			result.Err = err
		} else {
			result.Data = messages
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *MessageStore) GetByProviderAndMessageId(provider, providerMessageId string) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		var message model.Message
		if err := ms.db.Get(&message, `SELECT * FROM message
			WHERE provider = ? AND providerMessageId = ?`, provider, providerMessageId); err != nil {
			result.Err = err
		} else {
			result.Data = &message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *MessageStore) Search(params map[string]interface{}) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		query := "SELECT * FROM message"
		where := ""
		order := "DESC"
		args := []interface{}{}

		if since, ok := params["since"]; ok {
			where += " createdTime > ?"
			order = "ASC"
			args = append(args, since)
		}
		if until, ok := params["until"]; ok {
			where += sqlAndWhere(where)
			where += " createdTime < ?"
			args = append(args, until)
		}
		if to, ok := params["to"]; ok {
			where += sqlAndWhere(where)
			where += " toNumber = ?"
			args = append(args, to)
		}
		if status, ok := params["status"]; ok {
			where += sqlAndWhere(where)
			where += " status = ?"
			args = append(args, status)
		}
		if where != "" {
			query += " WHERE" + where
		}

		query += " ORDER BY createdTime " + order

		if limit, ok := params["limit"]; ok {
			query += " LIMIT ?"
			args = append(args, limit)
		}

		var messages []*model.Message
		if err := ms.db.Select(&messages, query, args...); err != nil {
			result.Err = err
		} else {
			length := len(messages)
			if order == "ASC" && length > 1 {
				// Reverse the messages slice
				for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
					messages[i], messages[j] = messages[j], messages[i]
				}
			}

			result.Data = messages
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *MessageStore) Save(message *model.Message) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		_, err := ms.db.Exec(`INSERT INTO message
			(
				id,
				toNumber,
				fromName,
				body,
				async,
				route,
				provider,
				providerMessageId,
				steps,
				status,
				createdTime,
				updatedTime
			)
			VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			message.Id,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Provider,
			message.ProviderMessageId,
			message.Steps,
			message.Status,
			message.CreatedTime,
			message.UpdatedTime,
		)
		if err != nil {
			result.Err = err
		} else {
			result.Data = message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *MessageStore) Update(message *model.Message) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		_, err := ms.db.Exec(`UPDATE message
			SET
				toNumber = ?,
				fromName = ?,
				body = ?,
				async = ?,
				route = ?,
				provider = ?,
				providerMessageId = ?,
				steps = ?,
				status = ?,
				createdTime = ?,
				updatedTime = ?
			WHERE id = ?`,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Provider,
			message.ProviderMessageId,
			message.Steps,
			message.Status,
			message.CreatedTime,
			message.UpdatedTime,
			message.Id,
		)
		if err != nil {
			result.Err = err
		} else {
			result.Data = message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/minchao/smsender/smsender/model"
)

const SqlMessageTable = `
CREATE TABLE IF NOT EXISTS message (
  id                varchar(40) COLLATE utf8_unicode_ci NOT NULL,
  toNumber          varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  fromName          varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  body              text COLLATE utf8_unicode_ci NOT NULL,
  async             tinyint(1) NOT NULL DEFAULT '0',
  route             varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  provider          varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  status            varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  originalMessageId varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  originalResponse  json DEFAULT NULL,
  originalReceipts  json DEFAULT NULL,
  createdTime       datetime(6) NOT NULL,
  updatedTime       datetime(6) DEFAULT NULL,
  sentTime          datetime(6) DEFAULT NULL,
  latency           int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY providerOriginalMessageId (provider, originalMessageId)
) DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`

type SqlMessageStore struct {
	*SqlStore
}

func NewSqlMessageStore(sqlStore *SqlStore) MessageStore {
	ms := &SqlMessageStore{sqlStore}

	ms.db.MustExec(SqlMessageTable)

	return ms
}

func (ms *SqlMessageStore) Get(id string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		var message model.MessageRecord
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

func (ms *SqlMessageStore) GetByIds(ids []string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		query, args, err := sqlx.In("SELECT * FROM message WHERE id IN (?)", ids)
		if err != nil {
			result.Err = err
			storeChannel <- result
			close(storeChannel)
			return
		}
		query = ms.db.Rebind(query)

		var messages []*model.MessageRecord
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

func (ms *SqlMessageStore) GetByProviderAndMessageId(provider, originalMessageId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		var message model.MessageRecord
		if err := ms.db.Get(&message, `SELECT * FROM message
			WHERE provider = ? AND originalMessageId = ?`, provider, originalMessageId); err != nil {
			result.Err = err
		} else {
			result.Data = &message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *SqlMessageStore) Search(params map[string]interface{}) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		query := "SELECT * FROM message"
		where := ""
		args := []interface{}{}

		if since, ok := params["since"]; ok {
			where += " createdTime > ?"
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

		query += " ORDER BY createdTime DESC"

		if limit, ok := params["limit"]; ok {
			query += " LIMIT ?"
			args = append(args, limit)
		}

		var messages []*model.MessageRecord
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

func (ms *SqlMessageStore) Save(message *model.MessageRecord) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		_, err := ms.db.Exec(`INSERT INTO message
			(
				id,
				toNumber,
				fromName,
				body,
				async,
				route,
				provider,
				status,
				originalMessageId,
				originalResponse,
				originalReceipts,
				createdTime,
				updatedTime,
				sentTime,
				latency
			)
			VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			message.Id,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Provider,
			message.Status,
			message.OriginalMessageId,
			message.OriginalResponse,
			message.OriginalReceipts,
			message.CreatedTime,
			message.UpdatedTime,
			message.SentTime,
			message.Latency,
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

func (ms *SqlMessageStore) Update(message *model.MessageRecord) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		_, err := ms.db.Exec(`UPDATE message
			SET
				toNumber = ?,
				fromName = ?,
				body = ?,
				async = ?,
				route = ?,
				provider = ?,
				status = ?,
				originalMessageId = ?,
				originalResponse = ?,
				originalReceipts = ?,
				createdTime = ?,
				updatedTime = ?,
				sentTime = ?,
				latency = ?
			WHERE id = ?`,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Provider,
			message.Status,
			message.OriginalMessageId,
			message.OriginalResponse,
			message.OriginalReceipts,
			message.CreatedTime,
			message.UpdatedTime,
			message.SentTime,
			message.Latency,
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

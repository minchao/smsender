package store

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/minchao/smsender/smsender/model"
)

const SqlMessageTable = `
CREATE TABLE IF NOT EXISTS message (
  id          varchar(40) COLLATE utf8_unicode_ci NOT NULL,
  toNumber    varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  fromName    varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  body        text COLLATE utf8_unicode_ci NOT NULL,
  async       tinyint(1) NOT NULL DEFAULT '0',
  route       varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  broker      varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  status      varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  original    json DEFAULT NULL,
  createdTime datetime(6) DEFAULT NULL,
  sentTime    datetime(6) DEFAULT NULL,
  latency     int(11) DEFAULT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci`

type SqlMessageStore struct {
	*SqlStore
}

func NewSqlMessageStore(sqlStore *SqlStore) MessageStore {
	ms := &SqlMessageStore{sqlStore}

	ms.db.MustExec(SqlMessageTable)

	return ms
}

func (ms *SqlMessageStore) Find(id string) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		var message model.Result
		if err := ms.db.Get(&message, `SELECT * FROM message WHERE id = ?`, id); err != nil {
			result.Err = err
		} else {
			if message.Original != nil {
				if original, err := unmarshalOriginal(message.Original); err == nil {
					message.Original = original
				}
			}
			result.Data = message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *SqlMessageStore) FindByIds(ids []string) StoreChannel {
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

		var messages []*model.Result
		if err := ms.db.Select(&messages, query, args...); err != nil {
			result.Err = err
		} else {
			for _, message := range messages {
				if message.Original != nil {
					if original, err := unmarshalOriginal(message.Original); err == nil {
						message.Original = original
					}
				}
			}
			result.Data = messages
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (ms *SqlMessageStore) Save(message *model.Result) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		var original *string
		if message.Original != nil {
			original, _ = marshalOriginal(message.Original)
		}

		_, err := ms.db.Exec(`INSERT INTO message
			(id, toNumber, fromName, body, async, route, broker, status, original, createdTime, sentTime, latency)
			VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			message.Id,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Broker,
			message.Status,
			original,
			message.CreatedTime,
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

func (ms *SqlMessageStore) Update(message *model.Result) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		var original *string
		if message.Original != nil {
			original, _ = marshalOriginal(message.Original)
		}

		_, err := ms.db.Exec(`UPDATE message
			SET
				toNumber = ?,
				fromName = ?,
				body = ?,
				async = ?,
				route = ?,
				broker = ?,
				status = ?,
				original = ?,
				createdTime = ?,
				sentTime = ?,
				latency = ?
			WHERE id = ?`,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Broker,
			message.Status,
			original,
			message.CreatedTime,
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

func marshalOriginal(original interface{}) (*string, error) {
	if result, err := json.Marshal(original); err != nil {
		return nil, err
	} else {
		str := string(result)
		return &str, nil
	}
}

func unmarshalOriginal(original interface{}) (interface{}, error) {
	var result interface{}
	if err := json.Unmarshal([]byte(original.([]uint8)), &result); err != nil {
		return nil, err
	}
	return result, nil
}

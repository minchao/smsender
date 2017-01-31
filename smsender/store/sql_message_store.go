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
  broker            varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  status            varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  originalMessageId varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  originalResponse  json DEFAULT NULL,
  originalReceipt   json DEFAULT NULL,
  createdTime       datetime(6) DEFAULT NULL,
  sentTime          datetime(6) DEFAULT NULL,
  receiptTime       datetime(6) DEFAULT NULL,
  latency           int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY brokerOriginalMessageId (broker, originalMessageId)
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
			result.Data = message
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
				broker,
				status,
				originalMessageId,
				originalResponse,
				originalReceipt,
				createdTime,
				sentTime,
				receiptTime,
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
			message.Broker,
			message.Status,
			message.OriginalMessageId,
			message.OriginalResponse,
			message.OriginalReceipt,
			message.CreatedTime,
			message.SentTime,
			message.ReceiptTime,
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
				broker = ?,
				status = ?,
				originalMessageId = ?,
				originalResponse = ?,
				originalReceipt = ?,
				createdTime = ?,
				sentTime = ?,
				receiptTime = ?,
				latency = ?
			WHERE id = ?`,
			message.To,
			message.From,
			message.Body,
			message.Async,
			message.Route,
			message.Broker,
			message.Status,
			message.OriginalMessageId,
			message.OriginalResponse,
			message.OriginalReceipt,
			message.CreatedTime,
			message.SentTime,
			message.ReceiptTime,
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

func (ms *SqlMessageStore) UpdateReceipt(receipt *model.MessageReceipt) StoreChannel {
	storeChannel := make(StoreChannel, 1)

	go func() {
		result := StoreResult{}

		_, err := ms.db.Exec(`UPDATE message
			SET
				status = ?,
				originalReceipt = ?,
				receiptTime = ?
			WHERE broker = ? AND originalMessageId = ?`,
			receipt.Status,
			model.MarshalJSON(receipt),
			receipt.CreatedTime,
			receipt.Broker,
			receipt.OriginalMessageId,
		)
		if err != nil {
			result.Err = err
		} else {
			result.Data = receipt
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

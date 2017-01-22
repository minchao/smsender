package model

import (
	"time"

	"github.com/rs/xid"
)

type StatusCode int

func (c StatusCode) String() string {
	return statusCodeMap[c]
}

const (
	// Default status, This should not be exported to client
	StatusInit StatusCode = iota
	// Received your API request to send a message
	StatusAccepted
	// The message is queued to be sent out
	StatusQueued
	// The message is in the process of dispatching to the upstream carrier
	StatusSending
	// The message was successfully accepted by the upstream carrie
	StatusSent
	// The message could not be sent to the upstream carrier
	StatusFailed
	// Received confirmation of message delivery from the upstream carrier
	StatusDelivered
	// Received that the message was not delivered from the upstream carrier
	StatusUndelivered
	// Received an undocumented status code from the upstream carrier
	StatusUnknown
)

var statusCodeMap = map[StatusCode]string{
	StatusInit:        "init",
	StatusAccepted:    "accepted",
	StatusQueued:      "queued",
	StatusSending:     "sending",
	StatusSent:        "sent",
	StatusFailed:      "failed",
	StatusDelivered:   "delivered",
	StatusUndelivered: "undelivered",
	StatusUnknown:     "unknown",
}

type Data struct {
	Id          string    `json:"id"`                 // Message Id
	To          string    `json:"to" db:"toNumber"`   // The destination phone number (E.164 format)
	From        string    `json:"from" db:"fromName"` // Sender Id (phone number or alphanumeric)
	Body        string    `json:"body"`               // The text of the message
	Async       bool      `json:"async,omitempty"`    // Enable a background sending mode that is optimized for bulk sending
	CreatedTime time.Time `json:"created_time" db:"createdTime"`
}

type MessageResult struct {
	Data
	SentTime          *time.Time  `json:"sent_time" db:"sentTime"`
	Latency           *int64      `json:"-"` // Millisecond
	Route             string      `json:"route"`
	Broker            string      `json:"broker"`
	Status            string      `json:"status"`
	OriginalMessageId *string     `json:"original_message_id" db:"originalMessageId"`
	OriginalResponse  interface{} `json:"original_response" db:"originalResponse"`
}

func NewMessageResult(message Message, broker string) *MessageResult {
	return &MessageResult{
		Data:   message.Data,
		Route:  message.Route,
		Broker: broker,
		Status: StatusSending.String(),
	}
}

func NewAsyncMessageResult(message Message) *MessageResult {
	result := MessageResult{
		Data:   message.Data,
		Route:  StatusUnknown.String(),
		Broker: StatusUnknown.String(),
		Status: StatusAccepted.String(),
	}
	if message.From == "" {
		result.From = StatusUnknown.String()
	}
	return &result
}

func NewMessageReceipt(originalMessageId, broker, status string, receipt interface{}, created time.Time) *MessageReceipt {
	return &MessageReceipt{
		OriginalMessageId: originalMessageId,
		Broker:            broker,
		Status:            status,
		OriginalReceipt:   receipt,
		CreatedTime:       created,
	}
}

type MessageReceipt struct {
	OriginalMessageId string      `json:"original_message_id"`
	Broker            string      `json:"broker"`
	Status            string      `json:"status"`
	OriginalReceipt   interface{} `json:"original_receipt" db:"originalReceipt"`
	CreatedTime       time.Time   `json:"created_time" db:"createdTime"`
}

type MessageRecord struct {
	MessageResult
	OriginalReceipt interface{} `json:"original_receipt" db:"originalReceipt"`
	ReceiptTime     *time.Time  `json:"receipt_time" db:"receiptTime"`
}

func NewMessageRecord(result MessageResult, receipt interface{}, receiptTime *time.Time) *MessageRecord {
	return &MessageRecord{
		MessageResult:   result,
		OriginalReceipt: receipt,
		ReceiptTime:     receiptTime,
	}
}

type Message struct {
	Data
	Route  string             `json:"route"`
	Result chan MessageResult `json:"-"`
}

func NewMessage(to, from, body string, async bool) *Message {
	message := Message{
		Data: Data{
			Id:          xid.New().String(),
			To:          to,
			From:        from,
			Body:        body,
			CreatedTime: time.Now(),
		},
		Route:  StatusInit.String(),
		Result: nil,
	}
	if async {
		message.Async = true
	} else {
		message.Result = make(chan MessageResult, 1)
	}
	return &message
}

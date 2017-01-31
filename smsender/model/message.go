package model

import (
	"time"

	"github.com/rs/xid"
)

type StatusCode int

func (c StatusCode) String() string {
	return statusCodeMap[c]
}

func statusStringToCode(status string) StatusCode {
	for k, v := range statusCodeMap {
		if v == status {
			return k
		}
	}
	return StatusUnknown
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
	// The message could not be sent to the upstream carrier
	StatusFailed
	// The message was successfully accepted by the upstream carrie
	StatusSent
	// Received an undocumented status code from the upstream carrier
	StatusUnknown
	// Received that the message was not delivered from the upstream carrier
	StatusUndelivered
	// Received confirmation of message delivery from the upstream carrier
	StatusDelivered
)

var statusCodeMap = map[StatusCode]string{
	StatusInit:        "init",
	StatusAccepted:    "accepted",
	StatusQueued:      "queued",
	StatusSending:     "sending",
	StatusFailed:      "failed",
	StatusSent:        "sent",
	StatusUnknown:     "unknown",
	StatusUndelivered: "undelivered",
	StatusDelivered:   "delivered",
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
	UpdatedTime       *time.Time `json:"updated_time" db:"updatedTime"`
	SentTime          *time.Time `json:"sent_time" db:"sentTime"`
	Latency           *int64     `json:"-"` // Millisecond
	Route             string     `json:"route"`
	Broker            string     `json:"broker"`
	Status            string     `json:"status"`
	OriginalMessageId *string    `json:"original_message_id" db:"originalMessageId"`
	OriginalResponse  JSON       `json:"original_response" db:"originalResponse"`
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
		Route:  "unknown",
		Broker: "unknown",
		Status: StatusAccepted.String(),
	}
	if message.From == "" {
		result.From = "unknown"
	}
	return &result
}

type MessageReceipt struct {
	OriginalMessageId string      `json:"-"`
	Broker            string      `json:"-"`
	Status            string      `json:"status"`
	OriginalReceipt   interface{} `json:"original_receipt"`
	CreatedTime       time.Time   `json:"created_time"`
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

type MessageRecord struct {
	MessageResult
	OriginalReceipts JSON `json:"original_receipts" db:"originalReceipts"`
}

func NewMessageRecord(result MessageResult, receipts JSON) *MessageRecord {
	return &MessageRecord{
		MessageResult:    result,
		OriginalReceipts: receipts,
	}
}

func (m *MessageRecord) GetReceipts() []MessageReceipt {
	var receipts []MessageReceipt
	m.OriginalReceipts.Unmarshal(&receipts)
	return receipts
}

func (m *MessageRecord) SetReceipts(receipts []MessageReceipt) {
	m.OriginalReceipts = MarshalJSON(receipts)
}

func (m *MessageRecord) AddReceipt(receipt MessageReceipt) {
	receipts := m.GetReceipts()
	receipts = append(receipts, receipt)
	m.SetReceipts(receipts)
}

func (m *MessageRecord) HandleReceipt(receipt MessageReceipt) {
	m.AddReceipt(receipt)

	if rStatus := statusStringToCode(receipt.Status); rStatus > StatusSent && rStatus > statusStringToCode(m.Status) {
		m.Status = receipt.Status
	}

	now := time.Now()
	m.UpdatedTime = &now
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
		Route:  "unknown",
		Result: nil,
	}
	if async {
		message.Async = true
	} else {
		message.Result = make(chan MessageResult, 1)
	}
	return &message
}

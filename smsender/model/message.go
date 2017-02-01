package model

import (
	"time"

	"github.com/rs/xid"
)

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
	Provider          string     `json:"provider"`
	Status            StatusCode `json:"status"`
	OriginalMessageId *string    `json:"original_message_id" db:"originalMessageId"`
	OriginalResponse  JSON       `json:"original_response" db:"originalResponse"`
}

func NewMessageResult(message Message, provider string) *MessageResult {
	return &MessageResult{
		Data:     message.Data,
		Route:    message.Route,
		Provider: provider,
		Status:   StatusSending,
	}
}

func NewAsyncMessageResult(message Message) *MessageResult {
	result := MessageResult{
		Data:     message.Data,
		Route:    "unknown",
		Provider: "unknown",
		Status:   StatusAccepted,
	}
	if message.From == "" {
		result.From = "unknown"
	}
	return &result
}

type MessageReceipt struct {
	OriginalMessageId string      `json:"-"`
	Provider          string      `json:"-"`
	Status            StatusCode  `json:"status"`
	OriginalReceipt   interface{} `json:"original_receipt"`
	CreatedTime       time.Time   `json:"created_time"`
}

func NewMessageReceipt(originalMessageId, provider string, status StatusCode, receipt interface{}, created time.Time) *MessageReceipt {
	return &MessageReceipt{
		OriginalMessageId: originalMessageId,
		Provider:          provider,
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

	if receipt.Status > StatusSent && receipt.Status > m.Status {
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

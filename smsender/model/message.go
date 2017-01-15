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
	StatusDelivered StatusCode = iota
	StatusFailed
	StatusSent
	StatusQueued
	StatusUnknown
)

var statusCodeMap = map[StatusCode]string{
	StatusDelivered: "delivered",
	StatusFailed:    "failed",
	StatusSent:      "sent",
	StatusQueued:    "queued",
	StatusUnknown:   "unknown",
}

type Data struct {
	Id          string    `json:"id"`
	To          string    `json:"to" db:"toNumber"`
	From        string    `json:"from" db:"fromName"`
	Body        string    `json:"body"`
	Async       bool      `json:"async,omitempty"`
	CreatedTime time.Time `json:"created_time" db:"createdTime"`
}

type Result struct {
	Data
	SentTime *time.Time  `json:"sent_time" db:"sentTime"`
	Latency  *int64      `json:"-"` // nanosecond
	Route    string      `json:"route"`
	Broker   string      `json:"broker"`
	Status   string      `json:"status"`
	Original interface{} `json:"original"`
}

type Message struct {
	Data
	Route  string      `json:"route"`
	Result chan Result `json:"-"`
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
		Route:  StatusUnknown.String(),
		Result: nil,
	}
	if async {
		message.Async = true
	} else {
		message.Result = make(chan Result, 1)
	}
	return &message
}

func NewResult(msg Message, broker string) *Result {
	return &Result{
		Data:   msg.Data,
		Route:  msg.Route,
		Broker: broker,
		Status: StatusUnknown.String(),
	}
}

func NewAsyncResult(msg Message) *Result {
	result := Result{
		Data:   msg.Data,
		Route:  StatusUnknown.String(),
		Broker: StatusUnknown.String(),
		Status: StatusQueued.String(),
	}
	if msg.From == "" {
		result.From = StatusUnknown.String()
	}
	return &result
}

package smsender

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
	Id   string `json:"id"`
	To   string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}

type Result struct {
	Data
	Route    string      `json:"route"`
	Broker   string      `json:"broker"`
	Status   string      `json:"status"`
	Original interface{} `json:"original"`
}

type Message struct {
	Data   Data
	Route  string
	Result chan Result
}

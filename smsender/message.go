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
	To   string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}

type Result struct {
	Id string `json:"id,omitempty"`
	Data
	Route     string `json:"route"`
	Broker    string `json:"broker"`
	Status    string `json:"status"`
	RawStatus string `json:"raw_status,omitempty"`
}

type Message struct {
	Data   Data
	Route  string
	Result chan Result
}

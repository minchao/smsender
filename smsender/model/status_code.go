package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type StatusCode int

func (c StatusCode) Value() (driver.Value, error) {
	return []byte(c.String()), nil
}

func (c *StatusCode) Scan(src interface{}) error {
	switch t := src.(type) {
	case []byte:
		code, err := statusStringToCode(string(t))
		if err != nil {
			return err
		}
		*c = code
	case nil:
		*c = StatusInit
	default:
		return errors.New("Incompatible type for StatusCode")
	}
	return nil
}

func (c StatusCode) String() string {
	return statusCodeMap[c]
}

func (c StatusCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *StatusCode) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("StatusCode: UnmarshalJSON on nil pointer")
	}
	var code string
	err := json.Unmarshal(data, &code)
	if err != nil {
		return err
	}
	statusCode, err := statusStringToCode(code)
	if err != nil {
		return err
	}
	*c = statusCode
	return nil
}

func statusStringToCode(status string) (StatusCode, error) {
	for k, v := range statusCodeMap {
		if v == status {
			return k, nil
		}
	}
	return 0, fmt.Errorf("StatusCode %s not exists", status)
}

// Status
const (
	StatusInit        StatusCode = iota // Default status, This should not be exported to client
	StatusAccepted                      // Received your API request to send a message
	StatusQueued                        // The message is queued to be sent out
	StatusSending                       // The message is in the process of dispatching to the upstream carrier
	StatusFailed                        // The message could not be sent to the upstream carrier
	StatusSent                          // The message was successfully accepted by the upstream carrie
	StatusUnknown                       // Received an undocumented status code from the upstream carrier
	StatusUndelivered                   // Received that the message was not delivered from the upstream carrier
	StatusDelivered                     // Received confirmation of message delivery from the upstream carrier
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

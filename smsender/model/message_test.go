package model

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewMessageRecord(t *testing.T) {
	ct, _ := time.Parse(time.RFC3339, "2017-01-01T00:00:03.1415926+08:00")
	originalMessageId := "b288anp82b37873aj510"

	j := new(bytes.Buffer)
	json.Compact(j, []byte(`
	{
      "id":"b288anp82b37873aj510",
      "to":"+886987654321",
      "from":"+1234567890",
      "body":"Happy New Year 2017",
      "created_time":"2017-01-01T00:00:03.1415926+08:00",
      "updated_time":"2017-01-01T00:00:03.1415926+08:00",
      "sent_time":"2017-01-01T00:00:03.1415926+08:00",
      "route":"dummy",
      "provider":"dummy",
      "status":"delivered",
      "original_message_id":"b288anp82b37873aj510",
      "original_response":{
        "response":"response"
      },
      "original_receipts":[
        {
          "status":"delivered",
          "original_receipt":{
            "message_id":"b288anp82b37873aj510",
            "status":"delivered"
          },
          "created_time":"2017-01-01T00:00:03.1415926+08:00"
        }
      ]
    }`))

	messageRecord := NewMessageRecord(
		MessageResult{
			Data: Data{
				Id:          "b288anp82b37873aj510",
				To:          "+886987654321",
				From:        "+1234567890",
				Body:        "Happy New Year 2017",
				Async:       false,
				CreatedTime: ct,
			},
			UpdatedTime:       &ct,
			SentTime:          &ct,
			Latency:           nil,
			Route:             "dummy",
			Provider:          "dummy",
			Status:            StatusDelivered,
			OriginalMessageId: &originalMessageId,
			OriginalResponse:  JSON(`{"response":"response"}`),
		},
		nil,
	)

	receipt := struct {
		MessageId string `json:"message_id"`
		Status    string `json:"status"`
	}{
		MessageId: originalMessageId,
		Status:    "delivered",
	}

	messageRecord.AddReceipt(*NewMessageReceipt(originalMessageId, "dummy", StatusDelivered, receipt, ct))

	record, err := json.Marshal(messageRecord)
	if err != nil {
		t.Error("MessageRecord marshal error:", err.Error())
	}
	if !reflect.DeepEqual(record, j.Bytes()) {
		t.Errorf("NewMessageRecord returned %s, want %s", record, j)
	}
}

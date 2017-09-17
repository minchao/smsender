package model

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewMessageRecord(t *testing.T) {
	dummy := "dummy"
	providerMessageID := "b288anp82b37873aj510"
	ct, _ := time.Parse(time.RFC3339, "2017-01-01T00:00:03.1415926+08:00")

	j := new(bytes.Buffer)
	json.Compact(j, []byte(`
	{
      "id":"b288anp82b37873aj510",
      "to":"+886987654321",
      "from":"+1234567890",
      "body":"Happy New Year 2017",
      "async":false,
      "route":"dummy",
      "provider":"dummy",
      "provider_message_id":"b288anp82b37873aj510",
      "steps":[
        {
          "stage":"platform",
          "data":null,
          "status":"accepted",
          "created_time":"2017-01-01T00:00:03.1415926+08:00"
        },
        {
          "stage":"queue",
          "data":null,
          "status":"sending",
          "created_time":"2017-01-01T00:00:03.1415926+08:00"
        },
        {
          "stage":"queue.response",
          "data":null,
          "status":"sent",
          "created_time":"2017-01-01T00:00:03.1415926+08:00"
        },
                {
          "stage":"carrier.receipt",
          "data":null,
          "status":"delivered",
          "created_time":"2017-01-01T00:00:03.1415926+08:00"
        }
      ],
      "status":"delivered",
      "created_time":"2017-01-01T00:00:03.1415926+08:00",
      "updated_time":"2017-01-01T00:00:03.1415926+08:00"
    }`))

	message := NewMessage("+886987654321", "+1234567890", "Happy New Year 2017", false)
	message.ID = "b288anp82b37873aj510"
	message.Route = &dummy
	message.Provider = &dummy
	message.CreatedTime = ct
	message.SetSteps([]MessageStep{
		{
			Stage:       StagePlatform,
			Status:      StatusAccepted,
			Data:        JSON{},
			CreatedTime: ct,
		},
	})

	step1 := NewMessageStepSending()
	step1.CreatedTime = ct
	message.HandleStep(step1)

	step2 := NewMessageResponse(StatusSent, nil, &providerMessageID)
	step2.CreatedTime = ct
	message.HandleStep(step2)

	step3 := NewMessageReceipt(providerMessageID, dummy, StatusDelivered, nil)
	step3.CreatedTime = ct
	message.HandleStep(step3)

	record, err := json.Marshal(message)
	if err != nil {
		t.Error("message marshal error:", err.Error())
	}
	if !reflect.DeepEqual(record, j.Bytes()) {
		t.Errorf("NewMessage returned %s, want %s", record, j)
	}
}

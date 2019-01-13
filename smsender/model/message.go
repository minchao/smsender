package model

import (
	"time"

	"github.com/rs/xid"
)

const (
	StagePlatform       = "platform"
	StageQueue          = "queue"
	StageQueueResponse  = "queue.response"
	StageCarrier        = "carrier"
	StageCarrierReceipt = "carrier.receipt"
)

type Message struct {
	ID                string     `json:"id"`                 // Message ID
	To                string     `json:"to" db:"toNumber"`   // The destination phone number (E.164 format)
	From              string     `json:"from" db:"fromName"` // Sender ID (phone number or alphanumeric)
	Body              string     `json:"body"`               // The text of the message
	Async             bool       `json:"async"`              // Enable a background sending mode that is optimized for bulk sending
	Route             *string    `json:"route"`
	Provider          *string    `json:"provider"`
	ProviderMessageID *string    `json:"provider_message_id" db:"providerMessageId"`
	Steps             JSON       `json:"steps"`
	Status            StatusCode `json:"status"`
	CreatedTime       time.Time  `json:"created_time" db:"createdTime"`
	UpdatedTime       *time.Time `json:"updated_time" db:"updatedTime"`
}

func NewMessage(to, from, body string, async bool) *Message {
	now := Now()
	message := Message{
		ID:          xid.New().String(),
		To:          to,
		From:        from,
		Body:        body,
		Async:       async,
		Status:      StatusAccepted,
		CreatedTime: now,
		UpdatedTime: &now,
	}
	message.AddStep(MessageStep{
		Stage:       StagePlatform,
		Status:      StatusAccepted,
		Data:        JSON{},
		CreatedTime: now,
	})
	return &message
}

func (m *Message) GetSteps() []MessageStep {
	var steps []MessageStep
	_ = m.Steps.Unmarshal(&steps)
	return steps
}

func (m *Message) SetSteps(steps []MessageStep) {
	m.Steps = MarshalJSON(steps)
}

func (m *Message) AddStep(step MessageStep) {
	steps := m.GetSteps()
	steps = append(steps, step)
	m.SetSteps(steps)
}

func (m *Message) HandleStep(wrap MessageStepWrap) {
	step := wrap.GetStep()
	m.AddStep(step)
	m.UpdatedTime = &step.CreatedTime

	switch step.Stage {
	case StageQueueResponse:
		m.ProviderMessageID = wrap.(*MessageResponse).ProviderMessageID
		m.Status = step.Status
	case StageCarrierReceipt:
		if step.Status > StatusSent && step.Status > m.Status {
			m.Status = step.Status
		}
	}
}

type MessageJob struct {
	Message
	Result chan Message
}

func NewMessageJob(to, from, body string, async bool) *MessageJob {
	job := MessageJob{
		Message: *NewMessage(to, from, body, async),
	}
	if !async {
		job.Result = make(chan Message, 1)
	}
	return &job
}

type MessageStep struct {
	Stage       string      `json:"stage"`
	Data        interface{} `json:"data"`
	Status      StatusCode  `json:"status"`
	CreatedTime time.Time   `json:"created_time" db:"createdTime"`
}

func (ms MessageStep) GetStep() MessageStep {
	return ms
}

func NewMessageStepSending() *MessageStep {
	return &MessageStep{
		Stage:       StageQueue,
		Data:        nil,
		Status:      StatusSending,
		CreatedTime: Now(),
	}
}

type MessageStepWrap interface {
	GetStep() MessageStep
}

type MessageResponse struct {
	MessageStep
	ProviderMessageID *string
}

func NewMessageResponse(status StatusCode, response interface{}, providerMessageID *string) *MessageResponse {
	return &MessageResponse{
		MessageStep: MessageStep{
			Stage:       StageQueueResponse,
			Data:        response,
			Status:      status,
			CreatedTime: Now(),
		},
		ProviderMessageID: providerMessageID,
	}
}

type MessageReceipt struct {
	MessageStep
	ProviderMessageID string
	Provider          string
}

func NewMessageReceipt(providerMessageID, provider string, status StatusCode, receipt interface{}) *MessageReceipt {
	return &MessageReceipt{
		MessageStep: MessageStep{
			Stage:       StageCarrierReceipt,
			Data:        receipt,
			Status:      status,
			CreatedTime: Now(),
		},
		ProviderMessageID: providerMessageID,
		Provider:          provider,
	}
}

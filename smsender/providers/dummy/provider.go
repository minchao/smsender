package dummy

import (
	"github.com/minchao/smsender/smsender/model"
)

type Provider struct {
	name string
}

func NewProvider(name string) *Provider {
	return &Provider{
		name: name,
	}
}

func (b Provider) Name() string {
	return b.name
}

func (b Provider) Send(message model.Message) *model.MessageResponse {
	return model.NewMessageResponse(model.StatusDelivered, nil, &message.Id)
}

func (b Provider) Callback(register func(webhook *model.Webhook), receiptsCh chan<- model.MessageReceipt) {
}

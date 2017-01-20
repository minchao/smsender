package dummy

import (
	"github.com/minchao/smsender/smsender/model"
)

type Broker struct {
	name string
}

func NewBroker(name string) *Broker {
	return &Broker{
		name: name,
	}
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(message *model.Message, result *model.MessageResult) {
	result.Status = model.StatusSent.String()
	result.OriginalMessageId = &result.Id
}

func (b Broker) Callback(webhooks *[]*model.Webhook, receiptsCh chan<- model.MessageReceipt) {}

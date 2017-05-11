package dummy

import (
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/plugin"
	"github.com/spf13/viper"
)

const name = "dummy"

func init() {
	plugin.RegisterProvider(name, Plugin)
}

func Plugin(config *viper.Viper) (model.Provider, error) {
	return New(name), nil
}

type Provider struct {
	name string
}

// New creates Dummy Provider.
func New(name string) *Provider {
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

package not_found

import "github.com/minchao/smsender/smsender/model"

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

func (b Provider) Send(message *model.Message, result *model.MessageResult) {
	result.Status = model.StatusFailed
	result.OriginalMessageId = &result.Id
	result.OriginalResponse = model.MarshalJSON(struct {
		Error string `json:"error"`
	}{
		Error: "no_route_matches",
	})
}

func (b Provider) Callback(register func(webhook *model.Webhook), receiptsCh chan<- model.MessageReceipt) {
}

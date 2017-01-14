package dummy

import "github.com/minchao/smsender/smsender/model"

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

func (b Broker) Send(msg *model.Message, result *model.Result) {
	result.Status = model.StatusSent.String()
}

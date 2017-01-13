package twilio

import (
	twilio "github.com/carlosdp/twiliogo"
	"github.com/minchao/smsender/smsender/model"
)

type Broker struct {
	name   string
	client *twilio.TwilioClient
}

type Config struct {
	Sid   string
	Token string
}

func (c Config) NewBroker(name string) *Broker {
	return &Broker{
		name:   name,
		client: twilio.NewClient(c.Sid, c.Token),
	}
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(msg *model.Message, result *model.Result) {
	resp, err := twilio.NewMessage(
		b.client,
		msg.From,
		msg.To,
		twilio.Body(msg.Body),
	)
	if err != nil {
		result.Status = model.StatusFailed.String()
		result.Original = err
	} else {
		result.Status = convertStatus(resp.Status).String()
		result.Original = resp
	}
}

func convertStatus(rawStatus string) model.StatusCode {
	switch rawStatus {
	case "delivered":
		return model.StatusDelivered
	case "failed", "undelivered":
		return model.StatusFailed
	case "sent":
		return model.StatusSent
	case "queued":
		return model.StatusQueued
	default:
		return model.StatusFailed
	}
}

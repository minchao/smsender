package twilio

import (
	twilio "github.com/carlosdp/twiliogo"
	"github.com/minchao/smsender/smsender"
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

func (b Broker) Send(msg *smsender.Message, result *smsender.Result) {
	resp, err := twilio.NewMessage(
		b.client,
		msg.From,
		msg.To,
		twilio.Body(msg.Body),
	)
	if err != nil {
		result.Status = smsender.StatusFailed.String()
		result.Original = err
	} else {
		result.Status = convertStatus(resp.Status).String()
		result.Original = resp
	}
}

func convertStatus(rawStatus string) smsender.StatusCode {
	switch rawStatus {
	case "delivered":
		return smsender.StatusDelivered
	case "failed", "undelivered":
		return smsender.StatusFailed
	case "sent":
		return smsender.StatusSent
	case "queued":
		return smsender.StatusQueued
	default:
		return smsender.StatusFailed
	}
}

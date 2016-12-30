package twilio

import (
	log "github.com/Sirupsen/logrus"
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

func (b Broker) Send(msg smsender.Message) {
	result := smsender.NewResult(msg, b)

	resp, err := twilio.NewMessage(
		b.client,
		msg.Data.From,
		msg.Data.To,
		twilio.Body(msg.Data.Body),
	)
	if err != nil {
		result.Status = smsender.StatusFailed.String()
		result.Original = err

		log.Errorf("broker '%s' send message failed: %v", b.Name(), err)
	} else {
		result.Status = convertStatus(resp.Status).String()
		result.Original = resp

		log.Infof("broker '%s' send message: %+v, %+v", b.Name(), msg, resp)
	}

	b.Result(msg.Result, *result)
}

func (b Broker) Result(c chan smsender.Result, r smsender.Result) {
	c <- r
	close(c)
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

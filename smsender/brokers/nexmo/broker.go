package nexmo

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"gopkg.in/njern/gonexmo.v1"
)

type Broker struct {
	name   string
	client *nexmo.Client
}

type Config struct {
	Key    string
	Secret string
}

func (c Config) NewBroker(name string) *Broker {
	client, err := nexmo.NewClientFromAPI(c.Key, c.Secret)
	if err != nil {
		log.Fatalf("Could not create the aws session: %s", err)
	}

	return &Broker{
		name:   name,
		client: client,
	}
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(msg *smsender.Message, result *smsender.Result) {
	message := &nexmo.SMSMessage{
		From: msg.From,
		To:   msg.To,
		Type: nexmo.Unicode,
		Text: msg.Body,
	}

	resp, err := b.client.SMS.Send(message)
	if err != nil {
		result.Status = smsender.StatusFailed.String()

		log.Errorf("broker '%s' send message failed: %v", b.Name(), err)
	} else {
		if resp.MessageCount > 0 {
			message := resp.Messages[0]

			result.Status = convertStatus(message.Status.String()).String()
			result.Original = resp
		}

		log.Infof("broker '%s' send message: %+v, %+v", b.Name(), msg, resp)
	}
}

func convertStatus(rawStatus string) smsender.StatusCode {
	switch rawStatus {
	case nexmo.ResponseSuccess.String():
		return smsender.StatusSent
	default:
		return smsender.StatusFailed
	}
}

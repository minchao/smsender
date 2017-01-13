package nexmo

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender/model"
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

func (b Broker) Send(msg *model.Message, result *model.Result) {
	message := &nexmo.SMSMessage{
		From: msg.From,
		To:   msg.To,
		Type: nexmo.Unicode,
		Text: msg.Body,
	}

	resp, err := b.client.SMS.Send(message)
	if err != nil {
		result.Status = model.StatusFailed.String()
		result.Original = model.BrokerError{Error: err.Error()}
	} else {
		if resp.MessageCount > 0 {
			respMsg := resp.Messages[0]

			result.Status = convertStatus(respMsg.Status.String()).String()
		} else {
			result.Status = model.StatusFailed.String()
		}
		result.Original = resp
	}
}

func convertStatus(rawStatus string) model.StatusCode {
	switch rawStatus {
	case nexmo.ResponseSuccess.String():
		return model.StatusSent
	default:
		return model.StatusFailed
	}
}

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

func (c Config) NewBroker(name string) Broker {
	client, err := nexmo.NewClientFromAPI(c.Key, c.Secret)
	if err != nil {
		log.Fatalf("Could not create the aws session: %s", err)
	}

	return Broker{
		name:   name,
		client: client,
	}
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(msg smsender.Message) {
	message := &nexmo.SMSMessage{
		From:  msg.Originator,
		To:    msg.Recipient,
		Type:  nexmo.Text,
		Text:  msg.Body,
		Class: nexmo.Standard,
	}

	resp, err := b.client.SMS.Send(message)
	if err != nil {
		log.Errorf("broker '%s' send message failed: %v", b.Name(), err)
	} else {
		log.Infof("broker '%s' send message: %+v, %+v", b.Name(), msg, resp)
	}
}

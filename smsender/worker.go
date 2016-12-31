package smsender

import (
	log "github.com/Sirupsen/logrus"
)

type worker struct {
	id     int
	sender *Sender
}

func (w worker) process(msg *Message) {
	log.Infof("worker '%d' process: %+v", w.id, msg)

	var (
		broker Broker
		result *Result
	)

	for _, r := range w.sender.routes {
		if r.Match(msg.To) {
			if msg.From == "" && r.From != "" {
				msg.From = r.From
			}
			msg.Route = r.Name
			broker = r.Broker
			break
		}
	}

	// No route matched, use the default broker
	if broker == nil {
		broker = w.sender.GetBroker(DefaultBroker)
	}

	result = NewResult(*msg, broker.Name())

	broker.Send(msg, result)

	if msg.Result != nil {
		msg.Result <- *result
	}
}

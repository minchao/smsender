package smsender

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type worker struct {
	id     int
	sender *Sender
}

func (w worker) process(msg *Message) {
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

	logger := log.WithFields(log.Fields{
		"message_id": msg.Id,
		"worker_id":  w.id,
		"broker":     broker.Name(),
	})
	logger.WithField("message", *msg).Info("worker process")

	result = NewResult(*msg, broker.Name())

	broker.Send(msg, result)

	sentTime := time.Now()
	result.SentTime = &sentTime

	logger = logger.WithField("latency", sentTime.Sub(msg.CreatedTime).Nanoseconds())

	switch result.Status {
	case StatusFailed.String():
		logger.WithField("result", *result).Error("broker send message failed")
	default:
		logger.WithField("result", *result).Info("broker send message")
	}

	if msg.Result != nil {
		msg.Result <- *result
	}
}

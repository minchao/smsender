package smsender

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender/model"
)

type worker struct {
	id     int
	sender *Sender
}

func (w worker) process(message *model.Message) {
	var (
		broker model.Broker
		result *model.Result
	)

	if match, ok := w.sender.Match(message.To); ok {
		if message.From == "" && match.From != "" {
			message.From = match.From
		}
		message.Route = match.Name
		broker = match.GetBroker()
	}

	// No route matched, use the default broker
	if broker == nil {
		broker = w.sender.GetBroker(DefaultBroker)
	}

	logger := log.WithFields(log.Fields{
		"message_id": message.Id,
		"worker_id":  w.id,
		"broker":     broker.Name(),
	})
	logger.WithField("message", *message).Info("worker process")

	result = model.NewResult(*message, broker.Name())

	broker.Send(message, result)

	sentTime := time.Now()
	result.SentTime = &sentTime

	logger = logger.WithField("latency", sentTime.Sub(message.CreatedTime).Nanoseconds())

	switch result.Status {
	case model.StatusFailed.String():
		logger.WithField("result", *result).Error("broker send message failed")
	default:
		logger.WithField("result", *result).Info("broker send message")
	}

	if message.Result != nil {
		message.Result <- *result
	}
}

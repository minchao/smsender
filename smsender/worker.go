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

	// Save the send record to db
	rchan1 := w.sender.store.Message().Save(result)

	broker.Send(message, result)

	sentTime := time.Now()
	latency := sentTime.Sub(message.CreatedTime).Nanoseconds()
	result.SentTime = &sentTime
	result.Latency = &latency

	rchan2 := w.sender.store.Message().Update(result)

	logger2 := logger.WithField("result", *result)
	switch result.Status {
	case model.StatusFailed.String():
		logger2.Error("broker send message failed")
	default:
		logger2.Info("broker send message")
	}

	if message.Result != nil {
		message.Result <- *result
	}

	if r := <-rchan1; r.Err != nil {
		logger.Errorf("store save error: %v", r.Err)
	}
	if r := <-rchan2; r.Err != nil {
		logger.Errorf("store update error: %v", r.Err)
	}
}

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
		result *model.MessageResult
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

	log1 := log.WithFields(log.Fields{
		"message_id": message.Id,
		"worker_id":  w.id,
		"broker":     broker.Name(),
	})
	log1.WithField("message", *message).Info("worker process")

	result = model.NewMessageResult(*message, broker.Name())

	// Save the send record to db
	rch := w.sender.store.Message().Save(result)

	broker.Send(message, result)

	sentTime := time.Now()
	latency := sentTime.Sub(message.CreatedTime).Nanoseconds() / int64(time.Millisecond)
	result.SentTime = &sentTime
	result.Latency = &latency

	log2 := log1.WithField("result", *result)
	switch result.Status {
	case model.StatusFailed.String():
		log2.Error("broker send message failed")
	default:
		log2.Info("broker send message")
	}

	if message.Result != nil {
		message.Result <- *result
	}

	if r := <-rch; r.Err != nil {
		log1.Errorf("store save error: %v", r.Err)
	}
	if r := <-w.sender.store.Message().Update(result); r.Err != nil {
		log1.Errorf("store update error: %v", r.Err)
	}
}

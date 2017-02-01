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
	rch := w.sender.store.Message().Save(model.NewMessageRecord(*result, nil))

	broker.Send(message, result)

	now := time.Now()
	latency := now.Sub(message.CreatedTime).Nanoseconds() / int64(time.Millisecond)
	result.UpdatedTime = &now
	result.SentTime = &now
	result.Latency = &latency

	log2 := log1.WithField("result", *result)
	switch result.Status {
	case model.StatusSent.String(), model.StatusDelivered.String():
		log2.Info("broker send message")
	case model.StatusFailed.String(), model.StatusUndelivered.String(), model.StatusUnknown.String():
		log2.Error("broker send message failed")
	default:
		// Unexpected status
		log2.Error("unexpected message status")
	}

	if message.Result != nil {
		message.Result <- *result
	}

	if r := <-rch; r.Err != nil {
		log1.Errorf("store save error: %v", r.Err)
		return
	}
	if r := <-w.sender.store.Message().Update(model.NewMessageRecord(*result, nil)); r.Err != nil {
		log1.Errorf("store update error: %v", r.Err)
	}
}

func (w worker) receipt(receipt model.MessageReceipt) {
	log1 := log.WithFields(log.Fields{
		"worker_id":           w.id,
		"original_message_id": receipt.OriginalMessageId,
	})
	log1.WithField("receipt", receipt).Info("handle the message receipt")

	r := <-w.sender.store.Message().GetByBrokerAndMessageId(receipt.Broker, receipt.OriginalMessageId)
	if r.Err != nil {
		log1.Errorf("receipt update error: message not found. %v", r.Err)
		return
	}

	message := r.Data.(*model.MessageRecord)
	message.HandleReceipt(receipt)

	if r := <-w.sender.store.Message().Update(message); r.Err != nil {
		log1.Errorf("receipt update error: %v", r.Err)
	}
}

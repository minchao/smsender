package smsender

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender/model"
)

type worker struct {
	id     int
	sender *Sender
}

func (w worker) process(job *model.MessageJob) {
	var (
		message  = job.Message
		provider = w.sender.Router.NotFoundProvider
	)

	if match, ok := w.sender.Router.Match(message.To); ok {
		if message.From == "" && match.From != "" {
			message.From = match.From
		}
		route := match.Name
		message.Route = &route
		provider = match.GetProvider()
	}

	p := provider.Name()
	message.Provider = &p

	log1 := log.WithFields(log.Fields{
		"message_id": message.Id,
		"worker_id":  w.id,
	})
	log1.WithField("message", message).Debug("worker process")

	// Save the send record to db
	rch := w.sender.store.Message().Save(&message)

	message.HandleStep(model.NewMessageStepSending())
	message.HandleStep(provider.Send(message))

	switch message.Status {
	case model.StatusSent, model.StatusDelivered:
		log1.Debug("successfully sent the message to the carrier")
	default:
		log1.WithField("message", message).Error("unable to send the message to the carrier")
	}

	if job.Result != nil {
		job.Result <- message
	}

	if r := <-rch; r.Err != nil {
		log1.Errorf("store save error: %v", r.Err)
		return
	}
	if r := <-w.sender.store.Message().Update(&message); r.Err != nil {
		log1.Errorf("store update error: %v", r.Err)
	}
}

func (w worker) receipt(receipt model.MessageReceipt) {
	log1 := log.WithFields(log.Fields{
		"worker_id":           w.id,
		"original_message_id": receipt.ProviderMessageId,
	})
	log1.WithField("receipt", receipt).Info("handle the message receipt")

	r := <-w.sender.store.Message().GetByProviderAndMessageId(receipt.Provider, receipt.ProviderMessageId)
	if r.Err != nil {
		log1.Errorf("receipt update error: message not found. %v", r.Err)
		return
	}

	message := r.Data.(*model.Message)
	message.HandleStep(receipt)

	if r := <-w.sender.store.Message().Update(message); r.Err != nil {
		log1.Errorf("receipt update error: %v", r.Err)
	}
}

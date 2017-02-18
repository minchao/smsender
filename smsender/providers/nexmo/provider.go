package nexmo

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/utils"
	"gopkg.in/njern/gonexmo.v1"
)

type Provider struct {
	name          string
	client        *nexmo.Client
	enableWebhook bool
	webhookPath   string
}

type Config struct {
	Key           string
	Secret        string
	EnableWebhook bool
}

func (c Config) NewProvider(name string) *Provider {
	client, err := nexmo.NewClientFromAPI(c.Key, c.Secret)
	if err != nil {
		log.Fatalf("Could not create the nexmo client: %s", err)
	}
	return &Provider{
		name:          name,
		client:        client,
		enableWebhook: c.EnableWebhook,
		webhookPath:   "/webhooks/" + name,
	}
}

func (b Provider) Name() string {
	return b.name
}

func (b Provider) Send(message model.Message) *model.MessageResponse {
	sms := &nexmo.SMSMessage{
		From: message.From,
		To:   message.To,
		Type: nexmo.Unicode,
		Text: message.Body,
	}

	resp, err := b.client.SMS.Send(sms)
	if err != nil {
		return model.NewMessageResponse(model.StatusFailed, model.ProviderError{Error: err.Error()}, nil)
	} else {
		var status model.StatusCode
		var providerMessageId *string
		if resp.MessageCount > 0 {
			respMsg := resp.Messages[0]

			status = convertStatus(respMsg.Status.String())
			providerMessageId = &respMsg.MessageID
		} else {
			status = model.StatusFailed
		}
		return model.NewMessageResponse(status, resp, providerMessageId)
	}
}

type DeliveryReceipt struct {
	Msisdn           string `json:"msisdn"`
	To               string `json:"to"`
	NetworkCode      string `json:"network-code"`
	MessageId        string `json:"messageId"`
	Price            string `json:"price"`
	Status           string `json:"status"`
	Scts             string `json:"scts"`
	ErrCode          string `json:"err-code"`
	MessageTimestamp string `json:"message-timestamp"`
}

// see https://docs.nexmo.com/messaging/sms-api/api-reference#delivery_receipt
func (b Provider) Callback(register func(webhook *model.Webhook), receiptsCh chan<- model.MessageReceipt) {
	if !b.enableWebhook {
		return
	}

	register(&model.Webhook{
		Path: b.webhookPath,
		Func: func(w http.ResponseWriter, r *http.Request) {
			var receipt DeliveryReceipt
			err := utils.GetInput(r.Body, &receipt, nil)
			if err != nil {
				log.Errorf("webhooks '%s' json unmarshal error: %+v", b.name, receipt)

				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if receipt.MessageId == "" || receipt.Status == "" {
				log.Infof("webhooks '%s' empty request body", b.name)

				// When you set the callback URL for delivery receipt,
				// Nexmo will send several requests to make sure that webhook was okay (status code 200).
				w.WriteHeader(http.StatusOK)
				return
			}

			receiptsCh <- *model.NewMessageReceipt(
				receipt.MessageId,
				b.Name(),
				convertDeliveryReceiptStatus(receipt.Status),
				receipt)

			w.WriteHeader(http.StatusOK)
		},
		Method: "POST",
	})
}

func convertStatus(rawStatus string) model.StatusCode {
	var status model.StatusCode
	switch rawStatus {
	case nexmo.ResponseSuccess.String():
		status = model.StatusSent
	default:
		status = model.StatusFailed
	}
	return status
}

func convertDeliveryReceiptStatus(rawStatus string) model.StatusCode {
	var status model.StatusCode
	switch rawStatus {
	case "accepted", "buffered":
		status = model.StatusSent
	case "delivered":
		status = model.StatusDelivered
	case "failed", "rejected":
		status = model.StatusUndelivered
	default:
		// expired, unknown
		status = model.StatusUnknown
	}
	return status
}

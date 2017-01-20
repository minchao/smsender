package twilio

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	twilio "github.com/carlosdp/twiliogo"
	"github.com/minchao/smsender/smsender/model"
)

type Broker struct {
	name          string
	client        *twilio.TwilioClient
	enableWebhook bool
	siteURL       string
	webhookURL    string
}

type Config struct {
	Sid           string
	Token         string
	EnableWebhook bool
	SiteURL       string
}

func (c Config) NewBroker(name string) *Broker {
	broker := &Broker{
		name:   name,
		client: twilio.NewClient(c.Sid, c.Token),
	}
	if c.EnableWebhook {
		if c.SiteURL == "" {
			log.Fatal("Could not create the twilio broker: SiteURL cannot be empty")
		}
		broker.enableWebhook = true
		broker.siteURL = c.SiteURL
		broker.webhookURL = c.SiteURL + "/webhooks/" + name
	}
	return broker
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(message *model.Message, result *model.MessageResult) {
	optionals := []twilio.Optional{twilio.Body(message.Body)}
	if b.enableWebhook {
		optionals = append(optionals, twilio.StatusCallback(b.webhookURL))
	}

	resp, err := twilio.NewMessage(
		b.client,
		message.From,
		message.To,
		optionals...,
	)
	if err != nil {
		result.Status = model.StatusFailed.String()
		result.OriginalResponse = err
	} else {
		result.Status = convertStatus(resp.Status)
		result.OriginalMessageId = &resp.Sid
		result.OriginalResponse = resp
	}
}

type DeliveryReceipt struct {
	MessageSid    string `json:"MessageSid"`
	ApiVersion    string `json:"ApiVersion"`
	From          string `json:"From"`
	To            string `json:"To"`
	AccountSid    string `json:"AccountSid"`
	SmsSid        string `json:"SmsSid"`
	SmsStatus     string `json:"SmsStatus"`
	MessageStatus string `json:"MessageStatus"`
}

// see https://www.twilio.com/docs/guides/sms/how-to-confirm-delivery
func (b Broker) Callback(webhooks *[]*model.Webhook, receiptsCh chan<- model.MessageReceipt) {
	if !b.enableWebhook {
		return
	}

	*webhooks = append(*webhooks, &model.Webhook{
		Path: "/webhooks/" + b.Name(),
		Func: func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			receipt := DeliveryReceipt{
				MessageSid:    r.Form.Get("MessageSid"),
				ApiVersion:    r.Form.Get("ApiVersion"),
				From:          r.Form.Get("From"),
				To:            r.Form.Get("To"),
				AccountSid:    r.Form.Get("AccountSid"),
				SmsSid:        r.Form.Get("SmsSid"),
				SmsStatus:     r.Form.Get("SmsStatus"),
				MessageStatus: r.Form.Get("MessageStatus"),
			}
			if receipt.MessageSid == "" || receipt.SmsStatus == "" {
				log.Infof("webhooks '%s' empty request body", b.name)

				w.WriteHeader(http.StatusBadRequest)
				return
			}

			receiptsCh <- *model.NewMessageReceipt(
				receipt.MessageSid,
				b.Name(),
				convertStatus(receipt.SmsStatus),
				receipt,
				time.Now())

			w.WriteHeader(http.StatusOK)
		},
		Method: "POST",
	})
}

func convertStatus(rawStatus string) string {
	var status model.StatusCode
	switch rawStatus {
	case "delivered":
		status = model.StatusDelivered
	case "failed", "undelivered":
		status = model.StatusFailed
	case "sent":
		status = model.StatusSent
	case "queued":
		status = model.StatusQueued
	default:
		status = model.StatusFailed
	}
	return status.String()
}

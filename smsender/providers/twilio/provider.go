package twilio

import (
	"errors"
	"net/http"

	twilio "github.com/carlosdp/twiliogo"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/plugin"
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

const name = "twilio"

func init() {
	plugin.RegisterProvider(name, Plugin)
}

func Plugin(c *config.Viper) (model.Provider, error) {
	return Config{
		Sid:           c.GetString("sid"),
		Token:         c.GetString("token"),
		EnableWebhook: c.GetBool("webhook.enable"),
		SiteURL:       config.GetString("http.siteURL"),
	}.New(name)
}

type Provider struct {
	name          string
	client        *twilio.TwilioClient
	enableWebhook bool
	siteURL       string
	webhookPath   string
}

type Config struct {
	Sid           string
	Token         string
	EnableWebhook bool
	SiteURL       string
}

// New creates Twilio Provider.
func (c Config) New(name string) (*Provider, error) {
	provider := &Provider{
		name:   name,
		client: twilio.NewClient(c.Sid, c.Token),
	}
	if c.EnableWebhook {
		if c.SiteURL == "" {
			return nil, errors.New("Could not create the twilio provider: SiteURL cannot be empty")
		}
		provider.enableWebhook = true
		provider.siteURL = c.SiteURL
		provider.webhookPath = "/webhooks/" + name
	}
	return provider, nil
}

func (b Provider) Name() string {
	return b.name
}

func (b Provider) Send(message model.Message) *model.MessageResponse {
	optionals := []twilio.Optional{twilio.Body(message.Body)}
	if b.enableWebhook {
		optionals = append(optionals, twilio.StatusCallback(b.siteURL+b.webhookPath))
	}

	resp, err := twilio.NewMessage(
		b.client,
		message.From,
		message.To,
		optionals...,
	)
	if err != nil {
		return model.NewMessageResponse(model.StatusFailed, err, nil)
	}

	return model.NewMessageResponse(convertStatus(resp.Status), resp, &resp.Sid)
}

type DeliveryReceipt struct {
	MessageSid    string `json:"MessageSid"`
	APIVersion    string `json:"ApiVersion"`
	From          string `json:"From"`
	To            string `json:"To"`
	AccountSid    string `json:"AccountSid"`
	SmsSid        string `json:"SmsSid"`
	SmsStatus     string `json:"SmsStatus"`
	MessageStatus string `json:"MessageStatus"`
}

// Callback see https://www.twilio.com/docs/guides/sms/how-to-confirm-delivery
func (b Provider) Callback(register func(webhook *model.Webhook), receiptsCh chan<- model.MessageReceipt) {
	if !b.enableWebhook {
		return
	}

	register(&model.Webhook{
		Path: b.webhookPath,
		Func: func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()

			receipt := DeliveryReceipt{
				MessageSid:    r.Form.Get("MessageSid"),
				APIVersion:    r.Form.Get("ApiVersion"),
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
				receipt)

			w.WriteHeader(http.StatusOK)
		},
		Method: "POST",
	})
}

func convertStatus(rawStatus string) model.StatusCode {
	var status model.StatusCode
	switch rawStatus {
	case "accepted", "queued", "sending", "sent":
		status = model.StatusSent
	case "delivered":
		status = model.StatusDelivered
	case "undelivered", "failed":
		status = model.StatusUndelivered
	default:
		status = model.StatusUnknown
	}
	return status
}

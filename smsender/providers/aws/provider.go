package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/plugin"
	"github.com/spf13/viper"
)

const name = "aws"

func init() {
	plugin.RegisterProvider(name, Plugin)
}

func Plugin(config *viper.Viper) (model.Provider, error) {
	return Config{
		Region: config.GetString("region"),
		ID:     config.GetString("id"),
		Secret: config.GetString("secret"),
	}.New(name)
}

type Provider struct {
	name string
	svc  *sns.SNS
}

type Config struct {
	Region string
	ID     string
	Secret string
}

// New creates AWS Provider.
func (c Config) New(name string) (*Provider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(
			c.ID,
			c.Secret,
			"",
		),
	})
	if err != nil {
		return nil, err
	}

	return &Provider{
		name: name,
		svc:  sns.New(sess),
	}, nil
}

func (b Provider) Name() string {
	return b.name
}

func (b Provider) Send(msg model.Message) *model.MessageResponse {
	req, resp := b.svc.PublishRequest(&sns.PublishInput{
		Message: aws.String(msg.Body),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Key": { // Required
				DataType:    aws.String("String"), // Required
				StringValue: aws.String("String"),
			},
		},
		PhoneNumber: aws.String(msg.To),
	})

	err := req.Send()

	if err != nil {
		return model.NewMessageResponse(model.StatusFailed, model.ProviderError{Error: err.Error()}, nil)
	}

	return model.NewMessageResponse(model.StatusSent, resp, resp.MessageId)
}

// Callback TODO: see http://docs.aws.amazon.com/sns/latest/dg/sms_stats_usage.html
func (b Provider) Callback(register func(webhook *model.Webhook), receiptsCh chan<- model.MessageReceipt) {
}

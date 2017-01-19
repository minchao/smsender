package aws

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/minchao/smsender/smsender/model"
)

type Broker struct {
	name string
	svc  *sns.SNS
}

type Config struct {
	Region string
	ID     string
	Secret string
}

func (c Config) NewBroker(name string) *Broker {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(
			c.ID,
			c.Secret,
			"",
		),
	})
	if err != nil {
		log.Fatalf("Could not create the aws session: %s", err)
	}

	return &Broker{
		name: name,
		svc:  sns.New(sess),
	}
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(msg *model.Message, result *model.Result) {
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
		result.Status = model.StatusFailed.String()
		result.OriginalResponse = model.BrokerError{Error: err.Error()}
	} else {
		result.Status = model.StatusSent.String()
		result.OriginalMessageId = resp.MessageId
		result.OriginalResponse = resp
	}
}

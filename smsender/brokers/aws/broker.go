package aws

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/minchao/smsender/smsender"
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

func (c Config) NewBroker(name string) Broker {
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

	return Broker{
		name: name,
		svc:  sns.New(sess),
	}
}

func (b Broker) Name() string {
	return b.name
}

func (b Broker) Send(msg smsender.Message) {
	req, resp := b.svc.PublishRequest(&sns.PublishInput{
		Message: aws.String(msg.Data.Body),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Key": { // Required
				DataType:    aws.String("String"), // Required
				StringValue: aws.String("String"),
			},
		},
		PhoneNumber: aws.String(msg.Data.To),
	})

	result := smsender.Result{
		Data:   msg.Data,
		Route:  msg.Route,
		Broker: b.Name(),
	}

	err := req.Send()

	if err != nil {
		result.Status = smsender.StatusFailed.String()

		log.Errorf("broker '%s' send message failed: %v", b.Name(), err)
	} else {
		result.Id = *resp.MessageId
		result.Status = smsender.StatusSent.String()

		log.Infof("broker '%s' send message: %+v, %+v", b.Name(), msg, resp)
	}

	b.Result(msg.Result, result)
}

func (b Broker) Result(c chan smsender.Result, r smsender.Result) {
	c <- r
	close(c)
}

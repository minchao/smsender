# SMSender

[![Build Status](https://travis-ci.org/minchao/smsender.svg?branch=master)](https://travis-ci.org/minchao/smsender)
[![Go Report Card](https://goreportcard.com/badge/github.com/minchao/smsender)](https://goreportcard.com/report/github.com/minchao/smsender)

A SMS server written in Go (Golang).

* Support various SMS brokers.
* Uses routes to determine which broker to send SMS.
* SMS delivery worker.
* RESTful API.

## Install

```
go get github.com/minchao/smsender
```

Using the [Glide](https://glide.sh/) to install dependency packages:

```
glide install
```

Creating a Configuration file:
 
```
cp ./config/config.default.yml ./config.yml
```

## Brokers

Support brokers

* [AWS SNS (SMS)](https://aws.amazon.com/sns/)
* [Nexmo](https://www.nexmo.com/)
* [Twilio](https://www.twilio.com/)

For example, registering a broker on the sender server.

Add the broker key and secret to config.yml:

```yaml
brokers:
  nexmo:
    key: "NEXMO_KEY"
    secret: "NEXMO_SECRET"
```

Add the following code to main.go:

```go
	nexmoBroker := nexmo.Config{
		Key:    config.GetString("brokers.nexmo.key"),
		Secret: config.GetString("brokers.nexmo.secret"),
	}.NewBroker("nexmo")
	
	sender := smsender.SMSender(config.GetInt("worker.num"))
	sender.AddBroker(nexmoBroker)
```

## Matching Routes

Route can be define a phone number pattern to be matched with broker.

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

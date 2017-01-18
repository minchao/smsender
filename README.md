# SMSender

[![Build Status](https://travis-ci.org/minchao/smsender.svg?branch=master)](https://travis-ci.org/minchao/smsender)
[![Go Report Card](https://goreportcard.com/badge/github.com/minchao/smsender)](https://goreportcard.com/report/github.com/minchao/smsender)

A SMS server written in Go (Golang).

* Support various SMS brokers.
* Uses routes to determine which broker to send SMS.
* SMS delivery worker.
* SMS delivery records.
* RESTful API.

## Requirements

* Go
* MySQL >= 5.7

## Installing

```bash
go get github.com/minchao/smsender
```

Using the [Glide](https://glide.sh/) to install dependency packages:

```bash
glide install
```

Creating a Configuration file:
 
```bash
cp ./config/config.default.yml ./config.yml
```

Setup the MySQL DSN:

```yaml
db:
  dsn: "user:password@tcp(localhost:3306)/dbname?parseTime=true&loc=Local"
```

Registering brokers on the sender server.

Add the broker key and secret to config.yml:

```yaml
brokers:
  nexmo:
    key: "NEXMO_KEY"
    secret: "NEXMO_SECRET"
```

Add the following code to main.go:

```go
    sender := smsender.SMSender(config.GetInt("worker.num"))
    
	nexmoBroker := nexmo.Config{
		Key:    config.GetString("brokers.nexmo.key"),
		Secret: config.GetString("brokers.nexmo.secret"),
	}.NewBroker("nexmo")
	
	sender.AddBroker(nexmoBroker)
	sender.Run()
```

Build:

```bash
go build -o bin/smsender
```

Run:

```bash
./bin/smsender
```

## Brokers

Support brokers

* [AWS SNS (SMS)](https://aws.amazon.com/sns/)
* [Nexmo](https://www.nexmo.com/)
* [Twilio](https://www.twilio.com/)

Need another broker? Just implement the [Broker](https://github.com/minchao/smsender/blob/master/smsender/model/broker.go) interface.

## Matching Routes

Route can be define a phone number pattern to be matched with broker.

## RESTful API

The API document is written in YAML and found in the [smsender-openapi.yaml](https://github.com/minchao/smsender/blob/master/smsender-openapi.yaml).
You can use the [Swagger Editor](http://editor.swagger.io/) to open the document.

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

# SMSender

[![Build Status](https://travis-ci.org/minchao/smsender.svg?branch=master)](https://travis-ci.org/minchao/smsender)
[![Go Report Card](https://goreportcard.com/badge/github.com/minchao/smsender)](https://goreportcard.com/report/github.com/minchao/smsender)

A SMS server written in Go (Golang).

* Support various SMS providers.
* Uses routes to determine which provider to send SMS.
* SMS delivery worker.
* SMS delivery records.
* SMS delivery receipt.
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

Registering providers on the sender server.

Add the provider key and secret to config.yml:

```yaml
providers:
  nexmo:
    key: "NEXMO_KEY"
    secret: "NEXMO_SECRET"
```

Add the following code to main.go:

```go
    sender := smsender.SMSender(config.GetInt("worker.num"))
    
	nexmoProvider := nexmo.Config{
		Key:    config.GetString("providers.nexmo.key"),
		Secret: config.GetString("providers.nexmo.secret"),
	}.NewProvider("nexmo")
	
	sender.AddProvider(nexmoProvider)
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

## Providers

Support providers

* [AWS SNS (SMS)](https://aws.amazon.com/sns/)
* [Nexmo](https://www.nexmo.com/)
* [Twilio](https://www.twilio.com/)

Need another provider? Just implement the [Provider](https://github.com/minchao/smsender/blob/master/smsender/model/provider.go) interface.

## Matching Routes

Route can be define a phone number pattern to be matched with provider.

## RESTful API

The API document is written in YAML and found in the [smsender-openapi.yaml](https://github.com/minchao/smsender/blob/master/smsender-openapi.yaml).
You can use the [Swagger Editor](http://editor.swagger.io/) to open the document.

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

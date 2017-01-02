# SMSender

A SMS server written in Go (Golang).

[![Build Status](https://travis-ci.org/minchao/smsender.svg?branch=master)](https://travis-ci.org/minchao/smsender)
[![Go Report Card](https://goreportcard.com/badge/github.com/minchao/smsender)](https://goreportcard.com/report/github.com/minchao/smsender)

## Install

```
go get github.com/minchao/smsender
```

Using the [Glide](https://glide.sh/) to install dependency packages

```
glide install
```

## Brokers

Support brokers

* [AWS SNS (SMS)](https://aws.amazon.com/sns/)
* [Nexmo](https://www.nexmo.com/)
* [Twilio](https://www.twilio.com/)

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/providers/aws"
	"github.com/minchao/smsender/smsender/providers/dummy"
	"github.com/minchao/smsender/smsender/providers/nexmo"
	"github.com/minchao/smsender/smsender/providers/twilio"
	config "github.com/spf13/viper"
)

func main() {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	sender := smsender.SMSender()

	sender.Router.AddProvider(dummy.NewProvider("dummy"))
	sender.Router.AddProvider(aws.Config{
		Region: config.GetString("providers.aws.region"),
		ID:     config.GetString("providers.aws.id"),
		Secret: config.GetString("providers.aws.secret"),
	}.NewProvider("aws"))
	sender.Router.AddProvider(nexmo.Config{
		Key:           config.GetString("providers.nexmo.key"),
		Secret:        config.GetString("providers.nexmo.secret"),
		EnableWebhook: config.GetBool("providers.nexmo.webhook.enable"),
	}.NewProvider("nexmo"))
	sender.Router.AddProvider(twilio.Config{
		Sid:           config.GetString("providers.twilio.sid"),
		Token:         config.GetString("providers.twilio.token"),
		EnableWebhook: config.GetBool("providers.twilio.webhook.enable"),
		SiteURL:       config.GetString("http.siteURL"),
	}.NewProvider("twilio"))
	sender.Router.LoadFromDB()

	api.InitAPI(sender)

	sender.Run()
}

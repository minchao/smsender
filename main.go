package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/providers/dummy"
	"github.com/minchao/smsender/smsender/web"
	config "github.com/spf13/viper"
)

func main() {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	sender := smsender.SMSender()

	provider := dummy.NewProvider("dummy")

	sender.Router.AddProvider(provider)
	sender.Router.LoadFromDB()

	api.InitAPI(sender)
	web.InitWeb(sender)

	sender.Run()
}

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/providers/dummy"
	config "github.com/spf13/viper"
)

func main() {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	sender := smsender.SMSender(config.GetInt("worker.num"))

	provider := dummy.NewProvider("dummy")

	sender.AddProvider(provider)
	sender.LoadRoutesFromDB()
	sender.InitWebhooks()
	go sender.Run()

	server := api.NewServer(sender)
	server.Run()
}

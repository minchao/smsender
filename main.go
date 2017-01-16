package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/brokers/dummy"
	config "github.com/spf13/viper"
)

func main() {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	sender := smsender.SMSender(config.GetInt("worker.num"))

	broker := dummy.NewBroker("dummy")

	sender.AddBroker(broker)
	go sender.Run()

	server := api.NewServer(sender)
	server.Run()
}

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	config "github.com/spf13/viper"
)

func main() {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	broker := smsender.NewDummyBroker("dummy")

	sender := smsender.SMSender(config.GetInt("worker.num"))
	sender.AddBroker(broker)
	sender.AddRoute("dummy", `.*`, broker.Name())
	go sender.Run()

	server := api.NewServer(config.GetString("api.addr"), sender)
	server.Run()
}

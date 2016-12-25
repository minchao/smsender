package main

import (
	"os"
	
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
)

func main() {

	broker := smsender.NewDummyBroker("dummy")

	sender := smsender.SMSender()
	sender.AddBroker(broker)
	sender.AddRoute("dummy", `.*`, broker.Name())
	go sender.Run()

	server := api.NewServer(os.Getenv("API_ADDR"), sender)
	server.Run()
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/providers/aws"
	"github.com/minchao/smsender/smsender/providers/dummy"
	"github.com/minchao/smsender/smsender/providers/nexmo"
	"github.com/minchao/smsender/smsender/providers/twilio"
	"github.com/minchao/smsender/smsender/web"
	config "github.com/spf13/viper"
)

func usage() {
	fmt.Println(`Usage: smsender [options]
Options are:
    -c, --config FILE  Configuration file path
    -h, --help         This help text`)
	os.Exit(0)
}

func handleSignals(s *smsender.Sender) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Infoln("Shutting down")
	go time.AfterFunc(60*time.Second, func() {
		os.Exit(1)
	})
	s.Shutdown()
	os.Exit(0)
}

func main() {
	var (
		help       bool
		configFile string
	)

	flag.BoolVar(&help, "h", false, "This help text")
	flag.BoolVar(&help, "help", false, "This help text")
	flag.StringVar(&configFile, "c", "", "Configuration file path")
	flag.StringVar(&configFile, "config", "", "Configuration file path")

	flag.Usage = usage
	flag.Parse()

	if help {
		usage()
	}

	if len(configFile) > 0 {
		config.SetConfigFile(configFile)
	} else {
		config.SetConfigName("config")
		config.AddConfigPath(".")
	}
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	log.Infof("Config path: %s", config.ConfigFileUsed())

	sender := smsender.NewSender()

	dummyProvider := dummy.NewProvider("dummy")
	sender.Router.AddProvider(dummyProvider)

	if ok := config.IsSet("providers.aws"); ok {
		provider := aws.Config{
			Region: config.GetString("providers.aws.region"),
			ID:     config.GetString("providers.aws.id"),
			Secret: config.GetString("providers.aws.secret"),
		}.NewProvider("aws")
		sender.Router.AddProvider(provider)
	}
	if ok := config.IsSet("providers.nexmo"); ok {
		provider := nexmo.Config{
			Key:           config.GetString("providers.nexmo.key"),
			Secret:        config.GetString("providers.nexmo.secret"),
			EnableWebhook: config.GetBool("providers.nexmo.webhook.enable"),
		}.NewProvider("nexmo")
		sender.Router.AddProvider(provider)
	}
	if ok := config.IsSet("providers.twilio"); ok {
		provider := twilio.Config{
			Sid:           config.GetString("providers.twilio.sid"),
			Token:         config.GetString("providers.twilio.token"),
			EnableWebhook: config.GetBool("providers.twilio.webhook.enable"),
			SiteURL:       config.GetString("http.siteURL"),
		}.NewProvider("twilio")
		sender.Router.AddProvider(provider)
	}

	sender.Router.LoadFromDB()

	api.InitAPI(sender)
	web.InitWeb(sender)

	go handleSignals(sender)

	sender.Run()
}

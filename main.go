package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/providers/dummy"
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

	sender := smsender.SMSender()

	provider := dummy.NewProvider("dummy")

	sender.Router.AddProvider(provider)
	sender.Router.LoadFromDB()

	api.InitAPI(sender)
	web.InitWeb(sender)

	sender.Run()
}

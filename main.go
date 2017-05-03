package main

import (
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
	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

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
	var configFile string
	var debug bool

	var rootCmd = &cobra.Command{
		Use:   "smsender",
		Short: "smsender",
		Long:  "A SMS server written in Go (Golang)",
		Run: func(cmd *cobra.Command, args []string) {

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

			if debug {
				log.SetLevel(log.DebugLevel)
				log.Debugln("Running in debug mode")
			}

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
		},
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	rootCmd.Execute()
}

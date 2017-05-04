package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/providers/aws"
	"github.com/minchao/smsender/smsender/providers/dummy"
	"github.com/minchao/smsender/smsender/providers/nexmo"
	"github.com/minchao/smsender/smsender/providers/twilio"
	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

func initEnv(cmd *cobra.Command) error {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if len(configFile) > 0 {
		config.SetConfigFile(configFile)
	} else {
		config.SetConfigName("config")
		config.AddConfigPath(".")
	}
	if err := config.ReadInConfig(); err != nil {
		return fmt.Errorf("Unable to read config file: %s", err)
	}

	log.Infof("Config path: %s", config.ConfigFileUsed())

	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Running in debug mode")
	}

	return nil
}

func initRouter(sender *smsender.Sender) error {
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

	return sender.Router.LoadFromDB()
}

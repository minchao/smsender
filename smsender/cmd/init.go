package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
		return fmt.Errorf("Unable to read config file %v", err)
	}

	log.Infof("Config path: %s", config.ConfigFileUsed())

	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Running in debug mode")
	}

	return nil
}

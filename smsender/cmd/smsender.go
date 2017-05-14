package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/api"
	"github.com/minchao/smsender/smsender/web"
	"github.com/spf13/cobra"
)

// Execute executes the smsender command.
func Execute() {
	var configFile string
	var debug bool

	var rootCmd = &cobra.Command{
		Use:   "smsender",
		Short: "smsender",
		Long:  "A SMS server written in Go (Golang)",
		Run:   rootCmdF,
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	rootCmd.AddCommand(sendCmd)

	rootCmd.Execute()
}

func rootCmdF(cmd *cobra.Command, args []string) {
	if err := initEnv(cmd); err != nil {
		log.Fatalln(err)
		return
	}

	sender := smsender.NewSender()
	api.InitAPI(sender)
	web.InitWeb(sender)

	go handleSignals(sender)

	sender.Run()
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

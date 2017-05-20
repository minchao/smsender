package cmd

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:     "routes",
	Short:   "List all routes",
	Example: `  routes`,
	RunE:    routesCmdF,
}

func routesCmdF(cmd *cobra.Command, args []string) error {
	if err := initEnv(cmd); err != nil {
		return err
	}

	sender := smsender.NewSender()
	resultJson, _ := json.MarshalIndent(sender.Router.GetAll(), "", "  ")

	log.Infof("Routes:\n%s", resultJson)

	return nil
}

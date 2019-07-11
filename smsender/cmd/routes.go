package cmd

import (
	"encoding/json"

	"github.com/minchao/smsender/smsender"
	log "github.com/sirupsen/logrus"
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
	resultJSON, _ := json.MarshalIndent(sender.Router.GetAll(), "", "  ")

	log.Infof("Routes:\n%s", resultJSON)

	return nil
}

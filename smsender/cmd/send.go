package cmd

import (
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/utils"
	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send message",
	Example: `  send --to +12345678900 --body "Hello, 世界"
  send --to +12345678900 --from smsender --body "Hello, 世界" --provider dummy`,
	RunE: sendCmdF,
}

func init() {
	sendCmd.Flags().StringP("to", "t", "", "The destination phone number (E.164 format)")
	sendCmd.Flags().StringP("from", "f", "", "Sender Id (phone number or alphanumeric)")
	sendCmd.Flags().StringP("body", "b", "", "The text of the message")
	sendCmd.Flags().StringP("provider", "p", "", "Provider name")
}

func sendCmdF(cmd *cobra.Command, args []string) error {
	if err := initEnv(cmd); err != nil {
		return err
	}

	to, err := cmd.Flags().GetString("to")
	if err != nil || to == "" {
		return errors.New("The to is required")
	}
	validate := validator.New()
	_ = validate.RegisterValidation("phone", utils.IsPhoneNumber)
	if err := validate.Var(to, "phone"); err != nil {
		return errors.New("Invalid phone number")
	}
	from, _ := cmd.Flags().GetString("from")
	body, err := cmd.Flags().GetString("body")
	if err != nil || body == "" {
		return errors.New("The body is required")
	}
	provider, _ := cmd.Flags().GetString("provider")

	config.Set("worker.num", 1)

	sender := smsender.NewSender()
	sender.InitWorkers()

	job := model.NewMessageJob(to, from, body, false)
	if provider != "" {
		job.Provider = &provider
	}

	queue := sender.GetMessagesChannel()
	queue <- job

	result := <-job.Result
	resultJSON, _ := json.MarshalIndent(result, "", "  ")

	log.Infof("Result:\n%s", resultJSON)

	return nil
}

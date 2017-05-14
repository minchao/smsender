package main

import (
	"github.com/minchao/smsender/smsender/cmd"

	// Register builtin stores.
	_ "github.com/minchao/smsender/smsender/store/sql"

	// Register builtin providers.
	_ "github.com/minchao/smsender/smsender/providers/aws"
	_ "github.com/minchao/smsender/smsender/providers/dummy"
	_ "github.com/minchao/smsender/smsender/providers/nexmo"
	_ "github.com/minchao/smsender/smsender/providers/twilio"
)

func main() {
	cmd.Execute()
}

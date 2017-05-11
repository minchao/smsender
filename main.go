package main

import (
	"github.com/minchao/smsender/smsender/cmd"
	_ "github.com/minchao/smsender/smsender/providers/aws"
	_ "github.com/minchao/smsender/smsender/providers/dummy"
	_ "github.com/minchao/smsender/smsender/providers/nexmo"
	_ "github.com/minchao/smsender/smsender/providers/twilio"
)

func main() {
	cmd.Execute()
}

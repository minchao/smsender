package plugin

import (
	"github.com/minchao/smsender/smsender/model"
	"github.com/spf13/viper"
)

func init() {
	ProviderFactories = map[string]ProviderFactory{}
}

type ProviderFactory func(config *viper.Viper) (model.Provider, error)

var ProviderFactories map[string]ProviderFactory

// RegisterProvider registers the specific ProviderFactory.
func RegisterProvider(name string, fn ProviderFactory) {
	ProviderFactories[name] = fn
}

package plugin

import (
	"github.com/minchao/smsender/smsender/model"
	"github.com/spf13/viper"
)

func init() {
	ProviderFactories = map[string]ProviderFactory{}
}

// ProviderFactory is a function that returns model.Provider implementation.
type ProviderFactory func(config *viper.Viper) (model.Provider, error)

// ProviderFactories is a map where provider name matches ProviderFactory.
var ProviderFactories map[string]ProviderFactory

// RegisterProvider registers the specific ProviderFactory.
func RegisterProvider(name string, fn ProviderFactory) {
	ProviderFactories[name] = fn
}

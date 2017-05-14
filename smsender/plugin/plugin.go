package plugin

import (
	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
	"github.com/spf13/viper"
)

func init() {
	StoreFactories = map[string]StoreFactory{}
	ProviderFactories = map[string]ProviderFactory{}
}

// StoreFactory is a function that returns store.Store implementation.
type StoreFactory func(config *viper.Viper) (store.Store, error)

// StoreFactories is a map where store name matches StoreFactory.
var StoreFactories map[string]StoreFactory

// RegisterStore registers the specific StoreFactory.
func RegisterStore(name string, fn StoreFactory) {
	StoreFactories[name] = fn
}

// ProviderFactory is a function that returns model.Provider implementation.
type ProviderFactory func(config *viper.Viper) (model.Provider, error)

// ProviderFactories is a map where provider name matches ProviderFactory.
var ProviderFactories map[string]ProviderFactory

// RegisterProvider registers the specific ProviderFactory.
func RegisterProvider(name string, fn ProviderFactory) {
	ProviderFactories[name] = fn
}

package providers

import (
	"fmt"
	"log"

	"github.com/cblomart/go-hue-bridge/config"
	"github.com/cblomart/go-hue-bridge/providers/domoticz"
	"github.com/cblomart/go-hue-bridge/providers/items"
)

var Providers map[string]Provider = make(map[string]Provider)

type Provider interface {
	GetLights() error
	On(id string) error
	Off(id string) error
}

func NewProvider(conf config.ProviderConfig) (Provider, error) {
	switch conf.Type {
	case "domoticz":
		d := domoticz.NewDomoticz(conf)
		return d, nil
	default:
		log.Fatalf("provider - create - unknown provider type: %s", conf.Type)
		return nil, fmt.Errorf("unkown prodicer type (%s)", conf.Type)
	}
}

func GetLights() {
	// loop over provider to scan for lights
	for name, provider := range Providers {
		// get light list
		err := provider.GetLights()
		if err != nil {
			log.Printf("provider - %s - cannot get lights from provider: %s", name, err)
			continue
		}
	}
	log.Printf("provider - %d lights present", len(items.Lights))
}

func GetProviderLights(id int) error {
	light := items.Lights[id]
	provider := Providers[light.Provider]
	return provider.GetLights()
}

func On(id int) error {
	light := items.Lights[id]
	provider := Providers[light.Provider]
	return provider.On(light.XID)
}

func Off(id int) error {
	light := items.Lights[id]
	provider := Providers[light.Provider]
	return provider.Off(light.XID)
}

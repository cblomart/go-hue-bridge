package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cblomart/go-hue-bridge/api"
	"github.com/cblomart/go-hue-bridge/config"
	"github.com/cblomart/go-hue-bridge/providers"
	"github.com/cblomart/go-hue-bridge/providers/items"
	"github.com/koron/go-ssdp"
	"gopkg.in/yaml.v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	file = flag.String("file", "/etc/go-hue-bridge.yaml", "Hue Bridge configuration")
)

func main() {
	// parse arguments
	flag.Parse()
	// print prog name
	log.Printf("Go Hue Bridge")
	// prepare config
	var hueConfig config.HueConfig
	// check config file presence
	if _, err := os.Stat(*file); errors.Is(err, os.ErrNotExist) {
		log.Printf("config file - %s - creating new config", *file)
		hueConfig = config.NewConfig()
		data, err := yaml.Marshal(&hueConfig)
		if err != nil {
			log.Fatalf("config file - %s - couldn't serialize: %s", *file, err)
		}
		err = ioutil.WriteFile(*file, data, 0666)
		if err != nil {
			log.Fatalf("config file - %s - couldn't write: %s", *file, err)
		}
	} else {
		// read config file
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			log.Fatalf("config file - %s - couldn't read: %s", *file, err)
		}
		err = yaml.Unmarshal(data, &hueConfig)
		if err != nil {
			log.Fatalf("config file - %s - couldn't read content: %s", *file, err)
		}
	}
	config.Config = hueConfig
	// initialize providers
	for _, provider := range config.Config.Providers {
		newprovider, err := providers.NewProvider(provider)
		if err == nil && newprovider != nil {
			providers.Providers[provider.Name] = newprovider
		}
	}
	providers.GetLights()
	log.Printf("hub - hub serial: %s", config.Config.Serial)
	log.Printf("hub - hub UUID: %s", config.Config.UUID)
	log.Printf("hub - hub ip address: %s", config.Config.IPAddress)
	log.Printf("hub - initial inventory: %d lights", len(items.Lights))

	// start ssdp advertisement
	_, err := ssdp.Advertise(
		"upnp:rootdevice", // send as "ST"
		"uuid:2fa00080-d000-11e1-9b23-001f80007bbe::upnp:rootdevice",       // send as "USN"
		fmt.Sprintf("http://%s:80/discovery.xml", config.Config.IPAddress), // send as "LOCATION"
		"FreeRTOS/6.0.5, UPnP/1.0, IpBridge/0.1",                           // send as "SERVER"
		100,                                                                // send as "maxAge" in "CACHE-CONTROL
	)
	if err != nil {
		log.Fatalf("hub - ssdp advertiser: %s", err)
	}
	log.Printf("hub - ssdp advertiser started")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	log.Printf("hub - route - /discovery.xml - main discovery endpoint")
	r.Get("/discovery.xml", api.Discovery)
	log.Printf("hub - route - /api - create users")
	r.Post("/api", api.Api)
	log.Printf("hub - route - /api/{userid}/lights - list lights")
	r.Get("/api/{userid}/lights", api.Lights)
	log.Printf("hub - route - /api/{userid}/lights/{id} - get light")
	r.Get("/api/{userid}/lights/{id}", api.LightInfo)
	log.Printf("hub - route - /api/{userid}/lights/{id}/state - set light")
	r.Put("/api/{userid}/lights/{id}/state", api.LightState)
	log.Printf("hub - http server starting")
	err = http.ListenAndServe(fmt.Sprintf("%s:80", hueConfig.IPAddress), r)
	if err != nil {
		log.Fatalf("hub - http server failed to start: %s", err)
	}
}

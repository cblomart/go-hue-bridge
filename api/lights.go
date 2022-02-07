package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cblomart/go-hue-bridge/providers"
	"github.com/cblomart/go-hue-bridge/providers/items"
)

type HueState struct {
	On        bool `json:"on"`
	Reachable bool `json:"reachable"`
}

type PointSymbol struct {
}

type HueLight struct {
	Name             string      `json:"name"`
	ManufacturerName string      `json:"manufacturername"`
	ModelID          string      `json:"modelid"`
	Version          string      `json:"swversion"`
	Type             string      `json:"type"`
	ID               string      `json:"uniqueid"`
	State            HueState    `json:"state"`
	PointSymbol      PointSymbol `json:"pointsymbol"`
}

func uniqueID(provider, name string) string {
	md := sha256.Sum256([]byte(fmt.Sprintf("%s/%s", provider, name)))
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x-%02x", md[1:2], md[2:3], md[4:5], md[6:7], md[8:9], md[10:11], md[12:13], md[14:15], md[16:17])
}

func GetLight(provider, name string, on bool) HueLight {
	// initialize dummy switch
	light := HueLight{
		Name:             name,
		ID:               uniqueID(provider, name),
		ManufacturerName: "OSRAM",
		ModelID:          "Plug 01",
		Type:             "On/off light",
		Version:          "v1.04.12",
		State: HueState{
			On:        on,
			Reachable: true,
		},
	}
	return light
}

func ToHueLight(l *items.Light) *HueLight {
	result := GetLight(l.Provider, l.Name, l.On)
	return &result
}

func ToHueLights(ls map[int]*items.Light) map[string]*HueLight {
	result := make(map[string]*HueLight)
	for i, l := range ls {
		result[strconv.Itoa(i)] = ToHueLight(l)
	}
	return result
}

func Lights(w http.ResponseWriter, r *http.Request) {
	// refresh providers
	providers.GetLights()
	// prep the response
	resp, err := json.Marshal(ToHueLights(items.Lights))
	if err != nil {
		log.Printf("http - lights - cannot convert to json: %s", err)
		w.WriteHeader(500)
		return
	}
	// send response
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte(resp))
	if err != nil {
		log.Fatalf("http - couldn't write light request: %s", err)
	}
}

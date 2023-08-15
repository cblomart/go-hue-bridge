package domoticz

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/cblomart/go-hue-bridge/config"
	"github.com/cblomart/go-hue-bridge/providers/items"
)

const (
	listlightsPath = "/json.htm?type=command&param=getdevices&filter=light&used=true&order=[Order]"
	setlightPath   = "/json.htm?type=command&param=switchlight&idx={id}&switchcmd={cmd}"
	getlightPath   = "/json.htm?type=command&param=getdevices&rid={id}"
)

type DomoticzLight struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	IDX    string `json:"idx"`
}

type DomoticzLights struct {
	Result []DomoticzLight `json:"result"`
}

type DomoticzStatus struct {
	Status  string `json:"status"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Domoticz struct {
	Name       string
	IPAddress  string
	Port       int
	Username   string
	Password   string
	SSL        bool
	StartIndex int
	client     *http.Client
}

func NewDomoticz(conf config.ProviderConfig) Domoticz {
	proto := "http"
	if conf.SSL {
		proto = "https"
	}
	auth := ""
	if len(conf.Username) > 0 && len(conf.Password) > 0 {
		auth = fmt.Sprintf("%s:*@", conf.Username)
	}
	log.Printf("domoticz - %s - new provider: %s://%s%s:%d", conf.Name, proto, auth, conf.IPAddress, conf.Port)
	return Domoticz{
		Name:       conf.Name,
		IPAddress:  conf.IPAddress,
		Port:       conf.Port,
		Username:   conf.Username,
		Password:   conf.Password,
		SSL:        conf.SSL,
		StartIndex: conf.StartIndex,
		client:     &http.Client{},
	}
}

func (d Domoticz) getURL(path string) string {
	proto := "http"
	if d.SSL {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s:%d%s", proto, d.IPAddress, d.Port, path)
}

func (d Domoticz) Set(id, state string) error {
	// prepare path
	path := strings.Replace(setlightPath, "{id}", id, -1)
	path = strings.Replace(path, "{cmd}", state, -1)
	// set lights from domoticz
	req, err := http.NewRequest("GET", d.getURL(path), nil)
	if err != nil {
		log.Printf("domoticz - %s - couldn't generate set request for %s to %s: %s", d.Name, id, state, err)
		return fmt.Errorf("couldn't generate set request for %s to %s", id, state)
	}
	resp, err := d.client.Do(req)
	if err != nil {
		log.Printf("domoticz - %s - couldn't set light %s to %s: %s", d.Name, id, state, err)
		return fmt.Errorf("couldn't set light %s to %s", id, state)
	}
	// read response
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	// convert to object
	var domoticzStatus DomoticzStatus
	err = json.Unmarshal(bodyBytes, &domoticzStatus)
	if err != nil {
		log.Printf("domoticz - %s - couldn't set light %s to %s: %s", d.Name, id, state, err)
		return fmt.Errorf("couldn't set light %s to %s", id, state)
	}
	if strings.EqualFold(domoticzStatus.Status, "error") {
		log.Printf("domoticz - %s - couldn't set light %s to %s: %s", d.Name, id, state, domoticzStatus.Message)
		return fmt.Errorf("couldn't set light %s to %s", id, state)
	}
	return nil
}

func (d Domoticz) On(id string) error {
	log.Printf("domoticz - %s - switching on %s", d.Name, id)
	return d.Set(id, "On")
}

func (d Domoticz) Off(id string) error {
	log.Printf("domoticz - %s - switching off %s", d.Name, id)
	return d.Set(id, "Off")
}

func (d Domoticz) GetLights() error {
	// get lights from domoticz
	req, err := http.NewRequest("GET", d.getURL(listlightsPath), nil)
	if err != nil {
		log.Printf("domoticz - %s - couldn't generate list request: %s", d.Name, err)
		return fmt.Errorf("couldn't generate list request")
	}
	resp, err := d.client.Do(req)
	if err != nil {
		log.Printf("domoticz - %s - couldn't get light list: %s", d.Name, err)
		return fmt.Errorf("couldn't get light list")
	}
	// read response
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	// convert to object
	var domoticzLights DomoticzLights
	err = json.Unmarshal(bodyBytes, &domoticzLights)
	if err != nil {
		log.Printf("domoticz - %s - couldn't read light list: %s", d.Name, err)
		return fmt.Errorf("couldn't read light list")
	}
	// convert domoticz lights to hue lights
	for i, l := range domoticzLights.Result {
		if _, ok := items.Lights[i+d.StartIndex]; ok {
			items.Lights[i+d.StartIndex].On = strings.EqualFold(l.Status, "On")
		} else {
			items.Lights[i+d.StartIndex] = &items.Light{Provider: d.Name, Name: l.Name, XID: l.IDX, On: strings.EqualFold(l.Status, "On")}
		}
	}
	log.Printf("domoticz - %s - returned %d lights", d.Name, len(domoticzLights.Result))
	return nil
}

func (d Domoticz) GetLight(id string) error {
	// prepare path
	path := strings.Replace(getlightPath, "{id}", id, -1)
	// set lights from domoticz
	req, err := http.NewRequest("GET", d.getURL(path), nil)
	if err != nil {
		log.Printf("domoticz - %s - couldn't generate get request for %s: %s", d.Name, id, err)
		return fmt.Errorf("couldn't generate get request for %s", id)
	}
	resp, err := d.client.Do(req)
	if err != nil {
		log.Printf("domoticz - %s - couldn't get light %s: %s", d.Name, id, err)
		return fmt.Errorf("couldn't get light %s", id)
	}
	// read response
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	// convert to object
	var domoticzLights DomoticzLights
	err = json.Unmarshal(bodyBytes, &domoticzLights)
	if err != nil {
		log.Printf("domoticz - %s - couldn't read light list: %s", d.Name, err)
		return fmt.Errorf("couldn't read light list")
	}
	if len(domoticzLights.Result) != 1 {
		log.Printf("domoticz - %s - light not found %s", d.Name, id)
		return nil
	}
	// find the light to update
	var lightItem *items.Light
	for _, l := range items.Lights {
		if l.XID == id {
			lightItem = l
			break
		}
	}
	if lightItem == nil {
		log.Printf("domoticz - %s - light not found %s", d.Name, id)
		return nil
	}
	lightItem.On = strings.EqualFold(domoticzLights.Result[0].Status, "On")
	log.Printf("domoticz - %s - updated light %s state", d.Name, id)
	return nil
}

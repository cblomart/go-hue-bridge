package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cblomart/go-hue-bridge/providers"
	"github.com/cblomart/go-hue-bridge/providers/items"
	"github.com/go-chi/chi/v5"
)

func LightInfo(w http.ResponseWriter, r *http.Request) {
	lightId := chi.URLParam(r, "id")
	lightId = strings.Replace(lightId, "\n", "", -1)
	lightId = strings.Replace(lightId, "\r", "", -1)
	id, err := strconv.ParseInt(lightId, 10, 0)
	if err != nil {
		log.Printf("http - light - cannot convert index (%s): %s", lightId, err)
		w.WriteHeader(500)
		return
	}
	// Refresh the provider
	err = providers.GetProviderLights(int(id))
	if err != nil {
		log.Printf("http - light - cannot refresh provider for %d: %s", id, err)
	}
	// prep the response
	resp, err := json.Marshal(ToHueLight(items.Lights[int(id)]))
	if err != nil {
		log.Printf("http - lights - cannot convert to json: %s", err)
		w.WriteHeader(500)
		return
	}
	// send response
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte(resp))
	if err != nil {
		log.Fatalf("http - couldn't write light info request: %s", err)
	}
}

func LightState(w http.ResponseWriter, r *http.Request) {
	// get light id
	lightId := chi.URLParam(r, "id")
	lightId = strings.Replace(lightId, "\n", "", -1)
	lightId = strings.Replace(lightId, "\r", "", -1)
	id, err := strconv.ParseInt(lightId, 10, 0)
	if err != nil {
		log.Printf("http - light - cannot convert index (%s): %s", lightId, err)
		w.WriteHeader(500)
		return
	}
	// get requested state
	req := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("http - light - cannot convert request for light (%d): %s", id, err)
		w.WriteHeader(500)
		return
	}
	for k := range req {
		reqState := fmt.Sprintf("%v", req[k])
		reqState = strings.Replace(reqState, "\n", "", -1)
		reqState = strings.Replace(reqState, "\r", "", -1)
		log.Printf("http - light - setting light %d '%s' to %v", id, k, reqState)
	}
	// set light state
	if req["on"].(bool) {
		err = providers.On(int(id))
	} else {
		err = providers.Off(int(id))
	}
	if err != nil {
		reqState := strconv.FormatBool(req["on"].(bool))
		log.Printf("http - light - cannot swith light (%d) '%s': %s", id, reqState, err)
		w.WriteHeader(500)
		return
	}
	// refresh provider
	err = providers.GetProviderLights(int(id))
	if err != nil {
		log.Printf("http - light - cannot get light '%d': %s", id, err)
		w.WriteHeader(500)
		return
	}
	// prep the response
	hueresp := make([]map[string]map[string]bool, 0)
	successResult := make(map[string]map[string]bool)
	successResult["success"] = make(map[string]bool)
	successResult["success"][fmt.Sprintf("/lights/%s/state/on", lightId)] = req["on"].(bool)
	hueresp = append(hueresp, successResult)
	resp, err := json.Marshal(hueresp)
	if err != nil {
		log.Printf("http - lights - cannot convert to json: %s", err)
		w.WriteHeader(500)
		return
	}
	// send response
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte(resp))
	if err != nil {
		log.Fatalf("http - couldn't write light switch request: %s", err)
	}
}

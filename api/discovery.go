package api

import (
	"net/http"
	"strings"

	"github.com/cblomart/go-hue-bridge/config"
)

const (
	discoveryResponse = `<root xmlns="urn:schemas-upnp-org:device-1-0">
<specVersion>
	<major>1</major>
	<minor>0</minor>
</specVersion>
<URLBase>http://{ip}:80/</URLBase>
<device>
	<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>
	<friendlyName>Philips hue ({ip})</friendlyName>
	<manufacturer>Royal Philips Electronics</manufacturer>
	<manufacturerURL>http://www.philips.com</manufacturerURL>
	<modelDescription>Philips hue Personal Wireless Lighting</modelDescription>
	<modelName>Philips hue bridge 2015</modelName>
	<modelNumber>BSB002</modelNumber>
	<modelURL>http://www.meethue.com</modelURL>
	<serialNumber>{serial}</serialNumber>
	<UDN>uuid:{uuid}</UDN>
</device>
</root>
`
)

func Discovery(w http.ResponseWriter, r *http.Request) {
	// prepare response
	resp := strings.ReplaceAll(discoveryResponse, "{ip}", config.Config.IPAddress)
	resp = strings.ReplaceAll(resp, "{serial}", config.Config.Serial)
	resp = strings.ReplaceAll(resp, "{uuid}", config.Config.UUID)
	resp = strings.ReplaceAll(resp, "\r", "")
	resp = strings.ReplaceAll(resp, "\n", "")
	// send response
	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(resp))
}

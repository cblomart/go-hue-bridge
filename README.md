# Golang Hue Bridge

The goal of this project is to build a simple golang hue brige implementing basic hue APIs.

It would mainly translate calls to the hue hub v1 api to calls to swith lights on or off.

It will support only on off: no colors or gradations.

## Hue Hub API

I don't have a fully docummented HUB API. It has been implemented by many other projects but sometimes fails with Alexa.

Discovery is done over SSDP responding to SSDP calls

```http
HTTP/1.1 200 OK
CACHE-CONTROL: max-age=100
EXT:
LOCATION: http://192.168.0.21:80/description.xml
SERVER: FreeRTOS/6.0.5, UPnP/1.0, IpBridge/0.1
ST: upnp:rootdevice
USN: uuid:2fa00080-d000-11e1-9b23-001f80007bbe::upnp:rootdevice
```

The discovery will provide a ```/discovery.xml``` url with details of the hub.

Request:
```http
GET /discovery.xml
```

Response:
```xml
<root xmlns="urn:schemas-upnp-org:device-1-0">
    <specVersion>
        <major>1</major>
        <minor>0</minor>
    </specVersion>
    <URLBase>http://{ip}:{port}/</URLBase>
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
```

Once Alexa has discovered the hub information it will use Hue API (v1) to discover lights.

The first thing Alexa will do is call the ``/api`` endpoint to associate with the hub and recieve a username.
The username will be leveraged in the uri of subsequent api calls to identify the hub. 
The username is a random hex string of 16bytes.

Request:
```http
GET /api
```

Response:
```json
[
    {
        "success": {
            "username": "{username}"
        }
    }
]
```

Once Alexa is associated to the hub it will query the lights endpoint to list all lights:
There i had to find the necessary attributes that needed to be presented. I chose to present OSRAM switches because these where the only on/off switches i could find on.
I started by adding all the attributes listed on ``Philips`` website but discovery didn't work. Alexa requested the lights endpoint ~10 times and failed to discover anything.
I found that some emulation softawre added the ``pointsymbol`` attribute. By adding it discovery did work.
Luckily i created a way to generate unique ids as described on ``Philips`` website as it turns out reverting to simple ids break the discovery too.

> *TODO*: check if exposing a specific device type is alright rather than copying existing ones

Request:
```http
GET /api/{username}/lights
```

Response:
```json
{
    "20": {
        "name":"Cave",
        "manufacturername":"OSRAM",
        "modelid":"Plug 01",
        "swversion":"v1.04.12",
        "type":"On/off light",
        "uniqueid":"00:7d:e8:a0:eb:37:00:f8-a4",
        "state": {
            "on":false,
            "effect":"none",
            "reachable":true,
            "alert":"none",
            "mode":"homeautomation"
        },
        "pointsymbol":{}
    },
    "21": {
        "name":"Cave - escalier",
        "manufacturername":"OSRAM",
        "modelid":"Plug 01",
        "swversion":"v1.04.12",
        "type":"On/off light",
        "uniqueid":"00:b5:64:46:fb:57:38:28-80",
        "state": {
            "on":false,
            "effect":"none",
            "reachable":true,
            "alert":"none",
            "mode":"homeautomation"
        },
        "pointsymbol":{}
    }
}
```
Once Alexa have discovered the lights list it will poll to each light to have their status.

Request:
```http
GET /api/{username}/lights/{id}
```

Response:
```json
{
        "name":"Cave",
        "manufacturername":"OSRAM",
        "modelid":"Plug 01",
        "swversion":"v1.04.12",
        "type":"On/off light",
        "uniqueid":"00:7d:e8:a0:eb:37:00:f8-a4",
        "state": {
            "on":false,
            "effect":"none",
            "xy":null,
            "reachable":true,
            "alert":"none",
            "mode":"homeautomation"
        },
        "pointsymbol":{}
    }
```

To change the state of a light Alexa will send a state request to a light with a put method

Request:
```http
PUT /api/{username}/lights/{id}/state
```

Request content:
```json
{
    "on":true
}
```

Response:
```json
[
    {
        "success": {
            "/lights/20/state/on": true
        }
    }
]
```
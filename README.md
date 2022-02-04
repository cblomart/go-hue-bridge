# Golang Hue Bridge

The goal of this project is to build a simple golang hue brige implementing basic hue APIs.

It would mainly translate calls to the hue hub v1 api to calls to swith lights on or off.

It will support only on off: no colors or gradations.

## Hue Hub API

I don't have a fully docummented HUB API. It has been implemented by many other projects but sometimes failes with Alexa.

### Discovery

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

From ```ha-bridge```.


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
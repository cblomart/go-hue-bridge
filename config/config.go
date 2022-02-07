package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"

	uuid "github.com/satori/go.uuid"
)

type Light struct {
	Name     string
	Provider string
	ID       string
}

type ProviderConfig struct {
	Type       string
	Name       string
	IPAddress  string
	Port       string
	Username   string
	Password   string
	SSL        bool
	StartIndex int
}

type HueConfig struct {
	Serial    string
	UUID      string
	IPAddress string
	Providers []ProviderConfig
}

var Config HueConfig

func randSerial() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func selectIp() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("config - detect ip: %s", err)
	}
	var ipAddrs []string
	for _, i := range ifaces {
		if i.Flags == net.FlagMulticast || i.Flags == net.FlagLoopback || i.Flags == net.FlagBroadcast || i.Flags == net.FlagPointToPoint {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatalf("config - detect ip: %s", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				continue
			}
			// process IP address
			ipAddrs = append(ipAddrs, ip.String())
		}
	}
	if len(ipAddrs) > 1 || len(ipAddrs) == 0 {
		log.Printf("config - detect ip: no ip detected or multiple, please adjust your config to use the proper one")
		for _, ip := range ipAddrs {
			log.Printf("config - available ip: %s", ip)
		}
		return ""
	}
	return ipAddrs[0]
}

func NewConfig() HueConfig {
	serial, _ := randSerial()
	uuid := uuid.NewV4()
	ipAddr := selectIp()
	return HueConfig{
		Serial:    serial,
		UUID:      uuid.String(),
		IPAddress: ipAddr,
	}
}

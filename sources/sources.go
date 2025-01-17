package sources

import (
	"fmt"
	"net"
)

func GetSource(name string, config map[string]string) (Source, error) {
	switch name {
	case "protonvpn":
		return NewProtonVPN()
	case "natpmp":
		gatewayIP := net.ParseIP(config["gatewayIP"])
		if gatewayIP == nil {
			return nil, fmt.Errorf("gatewayIP '%s' could not be parsed", config["gatewayIP"])
		}
		return NewNatPMP(gatewayIP)
	case "fake":
		return &Fake{}, nil
	default:
		return nil, fmt.Errorf("unknown source: %s", name)
	}
}

type Source interface {
	Get() (net.IP, int, error)
}

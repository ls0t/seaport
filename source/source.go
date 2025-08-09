package source

import (
	"fmt"
	"net"
	"time"
)

func Get(name string, options map[string]string) (Source, error) {
	switch name {
	case "protonvpn":
		return NewProtonVPN()
	case "natpmp":
		return NewNatPMP(options)
	case "gluetun":
		return NewGluetun(options)
	case "fake":
		return &Fake{}, nil
	default:
		return nil, fmt.Errorf("unknown source: %s", name)
	}
}

type Source interface {
	// Get returns the IP, port, and any error
	Get() (net.IP, int, error)

	// Refresh returns the interval on which refresh is initiated
	Refresh() time.Duration
}

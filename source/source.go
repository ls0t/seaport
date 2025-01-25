package source

import (
	"fmt"
	"net"
)

func Get(name string, options map[string]string) (Source, error) {
	switch name {
	case "protonvpn":
		return NewProtonVPN()
	case "natpmp":
		return NewNatPMP(options)
	case "fake":
		return &Fake{}, nil
	default:
		return nil, fmt.Errorf("unknown source: %s", name)
	}
}

type Source interface {
	Get() (net.IP, int, error)
}

package action

import (
	"context"
	"fmt"
	"net"
)

func Get(name string, options map[string]string) (Action, error) {
	switch name {
	case "qbittorrent":
		return NewQbittorrent(options), nil
	case "duckdns":
		return NewDuckDNS(options)
	case "freemyip":
		return NewFreeMyIP(options)
	default:
		return nil, fmt.Errorf("unknown action: %s", name)
	}
}

type Action interface {
	Act(ctx context.Context, ip net.IP, port int) error
	Name() string
}

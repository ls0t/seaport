package notify

import (
	"context"
	"fmt"
	"net"
)

func Get(name string, options map[string]string) (Notifier, error) {
	switch name {
	case "discord":
		return NewDiscord(options)
	default:
		return nil, fmt.Errorf("unknown notifier: %v", name)
	}
}

type Result struct {
	OldIP   net.IP
	OldPort int
	NewIP   net.IP
	NewPort int
	Err     error
}

type Notifier interface {
	Notify(ctx context.Context, result Result) error
}

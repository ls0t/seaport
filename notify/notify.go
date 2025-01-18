package notify

import (
	"context"
	"net"
)

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

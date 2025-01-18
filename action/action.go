package action

import (
	"context"
	"net"
)

type Action interface {
	Act(ctx context.Context, ip net.IP, port int) error
}

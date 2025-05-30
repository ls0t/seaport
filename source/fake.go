package source

import (
	"net"
	"time"
)

type Fake struct{}

func (f *Fake) Get() (net.IP, int, error) {
	return net.ParseIP("1.2.3.4"), 1234, nil
}

func (f *Fake) Refresh() time.Duration {
	return time.Hour
}

package action

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
)

type Transmission struct {
	url      string
	username string
	password string
}

func NewTransmission(options map[string]string) (Action, error) {
	url := strings.TrimSpace(options["url"])
	if url == "" {
		url = "http://localhost:9091"
	}
	username := strings.TrimSpace(options["username"])
	password := strings.TrimSpace(options["password"])

	return &Transmission{
		url:      url,
		username: username,
		password: password,
	}, nil
}

func (t *Transmission) Act(ctx context.Context, ip net.IP, port int) error {
	client, err := transmission.New(t.url)
	if err != nil {
		return fmt.Errorf("creating transmission client: %v", err)
	}

	useNatpmp := false
	err = client.SetSession(ctx, &transmission.SetSessionReq{
		PeerPort:              &port,
		PortForwardingEnabled: &useNatpmp,
	})
	if err != nil {
		return fmt.Errorf("when setting port: %v", err)
	}

	return nil
}

func (t *Transmission) Name() string {
	return "transmission"
}

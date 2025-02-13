package action

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/autobrr/go-qbittorrent"
)

type Qbittorrent struct {
	client *qbittorrent.Client
}

func NewQbittorrent(options map[string]string) Action {
	url := strings.TrimSpace(options["url"])
	if url == "" {
		url = "http://localhost:8080"
	}
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     url,
		Username: strings.TrimSpace(options["username"]),
		Password: strings.TrimSpace(options["password"]),
	})
	return &Qbittorrent{
		client: client,
	}
}

func (q *Qbittorrent) Act(ctx context.Context, ip net.IP, port int) error {
	err := q.client.LoginCtx(ctx)
	if err != nil {
		return fmt.Errorf("could not log into client: %w", err)
	}

	m := map[string]any{"listen_port": port, "random_port": "false", "upnp": "false"}
	err = q.client.SetPreferencesCtx(ctx, m)
	if err != nil {
		return fmt.Errorf("when updating qbittorrent settings: %w", err)
	}
	return nil
}

func (q *Qbittorrent) Name() string {
	return "qbittorrent"
}

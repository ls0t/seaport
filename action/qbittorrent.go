package action

import (
	"context"
	"fmt"
	"net"

	"github.com/autobrr/go-qbittorrent"
)

type QbittorrentConfig struct {
	Host     string
	Username string
	Password string
}

type Qbittorrent struct {
	client *qbittorrent.Client
}

func NewQbittorrent(options map[string]string) Action {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     options["host"],
		Username: options["username"],
		Password: options["password"],
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

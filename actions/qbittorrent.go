package actions

import (
	"context"
	"fmt"
	"net"

	"github.com/autobrr/go-qbittorrent"
)

/*
	host:     "http://localhost:8080",
	username: "admin",
	password: "adminadmin",
*/

type QbittorrentConfig struct {
	Host     string
	Username string
	Password string
}

type Qbittorrent struct {
	client *qbittorrent.Client
}

func NewQbittorrent(config QbittorrentConfig) Action {
	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     config.Host,
		Username: config.Username,
		Password: config.Password,
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

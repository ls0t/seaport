package notify

import (
	"context"
	"os"
	"strings"

	"github.com/disgoorg/disgo/discord"
	wh "github.com/disgoorg/disgo/webhook"
)

type Discord struct {
	client wh.Client
}

func NewDiscord(options map[string]string) (Notifier, error) {
	webhook := options["webhook"]
	if strings.TrimSpace(webhook) == "" {
		webhook = os.Getenv("SEAPORT_WEBHOOK_URL")
	}

	client, err := wh.NewWithURL(webhook)
	if err != nil {
		return nil, err
	}

	return &Discord{
		client: client,
	}, nil
}

func (d *Discord) Notify(ctx context.Context, result Result) error {
	var msg discord.WebhookMessageCreate
	if result.OldIP == nil {
		msg = discord.NewWebhookMessageCreateBuilder().
			SetContentf("Updated to %s:%d", result.NewIP, result.NewPort).
			Build()
	} else {
		msg = discord.NewWebhookMessageCreateBuilder().
			SetContentf("Updated from %s:%d to %s:%d", result.OldIP, result.OldPort, result.NewIP, result.NewPort).
			Build()
	}

	_, err := d.client.CreateMessage(msg)
	return err
}

package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ls0t/seeport/actions"
	"github.com/ls0t/seeport/notifiers"
	"github.com/ls0t/seeport/sources"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer stop()

	var ip net.IP
	port := 0
	var err error
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	source, err := sources.GetSource("protonvpn", nil)
	if err != nil {
		log.Fatalf("creating source failed: %v", err)
	}

	qbit := actions.NewQbittorrent(actions.QbittorrentConfig{
		Host:     "http://localhost:8080",
		Username: "admin",
		Password: "adminadmin",
	})
	actions := []actions.Action{qbit}

	discord, err := notifiers.NewDiscord(os.Getenv("SEEPORT_DISCORD_WEBHOOK"))
	if err != nil {
		log.Fatalf("creating notifier failed: %v", err)
	}
	notifiers := []notifiers.Notifier{discord}

	for {
		ip, port = tick(ctx, source, actions, notifiers, ip, port)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Print("exiting")
			return
		}
	}
}

func tick(ctx context.Context, source sources.Source, actions []actions.Action, toNotify []notifiers.Notifier, oldIP net.IP, oldPort int) (net.IP, int) {
	newIP, newPort, err := source.Get()
	if err != nil {
		log.Printf("error: %v", err)
		return oldIP, oldPort
	}
	if newPort != oldPort {
		log.Printf("updating from port %d to %d", oldPort, newPort)
		var results []notifiers.Result
		for _, action := range actions {
			err = action.Act(ctx, newIP, newPort)
			results = append(results, notifiers.Result{
				OldIP:   oldIP,
				OldPort: oldPort,
				NewIP:   newIP,
				NewPort: newPort,
				Err:     err,
			})
			if err != nil {
				log.Printf("error performing action: %v", err)
			}
		}

		log.Printf("updated to %s:%d", newIP, newPort)
		for _, notifier := range toNotify {
			for _, result := range results {
				err = notifier.Notify(ctx, result)
				if err != nil {
					log.Printf("error notifying: %v", err)
				}
			}
		}
	}

	return newIP, newPort
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ls0t/seeport/actions"
	"github.com/ls0t/seeport/sources"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer stop()

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

	for {
		port = tick(ctx, source, actions, port)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Print("exiting")
			return
		}
	}
}

func tick(ctx context.Context, source sources.Source, actions []actions.Action, oldPort int) int {
	externalIP, newPort, err := source.Get()
	if err != nil {
		log.Printf("error: %v", err)
		return oldPort
	}
	if newPort != oldPort {
		log.Printf("updating from port %d to %d", oldPort, newPort)
		for _, action := range actions {
			err = action.Act(ctx, externalIP, newPort)
			if err != nil {
				log.Printf("error updating qbittorrent: %v", err)
				return oldPort
			}
		}

		oldPort = newPort
		log.Printf("updated qbittorrent to %s:%d", externalIP, newPort)
	}
	return oldPort
}

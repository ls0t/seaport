package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/autobrr/go-qbittorrent"
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

	for {
		port = tick(ctx, source, port)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Print("exiting")
			return
		}
	}
}

func tick(ctx context.Context, source sources.Source, oldPort int) int {
	externalIP, newPort, err := source.Get()
	if err != nil {
		log.Printf("error: %v", err)
		return oldPort
	}
	if newPort != oldPort {
		log.Printf("updating from port %d to %d", oldPort, newPort)
		err = updateQbit(ctx, newPort)
		if err != nil {
			log.Printf("error updating qbittorrent: %v", err)
			return oldPort
		}
		oldPort = newPort
		log.Printf("updated qbittorrent to %s:%d", externalIP, newPort)
	}
	return oldPort
}

func updateQbit(ctx context.Context, externalPort int) error {
	qbitClient := qbittorrent.NewClient(qbittorrent.Config{
		Host:     "http://localhost:8080",
		Username: "admin",
		Password: "adminadmin",
	})

	err := qbitClient.LoginCtx(ctx)
	if err != nil {
		return fmt.Errorf("could not log into client: %w", err)
	}

	m := map[string]any{"listen_port": externalPort, "random_port": "false", "upnp": "false"}
	err = qbitClient.SetPreferencesCtx(ctx, m)
	if err != nil {
		return fmt.Errorf("couldn't update qbittorrent settings: %w", err)
	}
	return nil
}

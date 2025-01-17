package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/autobrr/go-qbittorrent"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer stop()

	gatewayIP := net.ParseIP("10.2.0.1")
	port := 0
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		port = tick(ctx, gatewayIP, port)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Print("exiting")
			return
		}
	}
}

func tick(ctx context.Context, gatewayIP net.IP, oldPort int) int {
	externalIP, newPort, err := getExternalPort(gatewayIP)
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

func getExternalPort(gatewayIP net.IP) (net.IP, int, error) {
	client := natpmp.NewClient(gatewayIP)

	portResponse, err := client.AddPortMapping("tcp", 0, 1, 60)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting tcp port mapping: %w", err)
	}
	externalPort := int(portResponse.MappedExternalPort)

	portResponse, err = client.AddPortMapping("udp", 0, externalPort, 60)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting udp port mapping: %w", err)
	}
	externalPort = int(portResponse.MappedExternalPort)

	addressResponse, err := client.GetExternalAddress()
	if err != nil {
		return nil, 0, fmt.Errorf("error getting external address: %w", err)
	}

	ip := net.IP(addressResponse.ExternalIPAddress[:])
	return ip, externalPort, nil
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

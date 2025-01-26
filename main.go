package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ls0t/seaport/action"
	"github.com/ls0t/seaport/config"
	"github.com/ls0t/seaport/notify"
	"github.com/ls0t/seaport/source"
)

var (
	configFilenameArg string
)

func init() {
	flag.StringVar(&configFilenameArg, "config", "seaport.yaml", "yaml config file")
}

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer stop()

	var ip net.IP
	port := 0
	var err error
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	f, err := os.Open(configFilenameArg)
	if err != nil {
		log.Fatalf("reading config file failed: %v", err)
	}
	defer f.Close()

	c, err := config.FromReader(f)
	if err != nil {
		log.Fatalf("parsing config failed: %v", err)
	}

	source, err := source.Get(c.Source.Name, c.Source.Options)
	if err != nil {
		log.Fatalf("creating source failed: %v", err)
	}

	actions := []action.Action{}
	for _, act := range c.Actions {
		a, err := action.Get(act.Name, act.Options)
		if err != nil {
			log.Fatalf("creating action failed: %v", err)
		}
		actions = append(actions, a)
	}

	notifiers := []notify.Notifier{}
	for _, not := range c.Notifiers {
		n, err := notify.Get(not.Name, not.Options)
		if err != nil {
			log.Fatalf("creating notifier failed: %v", err)
		}
		notifiers = append(notifiers, n)
	}

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

func tick(ctx context.Context, source source.Source, actions []action.Action, notifiers []notify.Notifier, oldIP net.IP, oldPort int) (net.IP, int) {
	newIP, newPort, err := source.Get()
	if err != nil {
		log.Printf("error: %v", err)
		return oldIP, oldPort
	}
	if newPort != oldPort {
		log.Printf("updating from port %d to %d", oldPort, newPort)
		var results []notify.Result
		for _, action := range actions {
			err = action.Act(ctx, newIP, newPort)
			results = append(results, notify.Result{
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
		for _, notifier := range notifiers {
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

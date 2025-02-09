package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
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
	displayVersionArg bool
	debugArg          bool
	version           = "dev" // populated from build flags
)

func init() {
	flag.StringVar(&configFilenameArg, "config", "seaport.yaml", "yaml config file")
	flag.BoolVar(&displayVersionArg, "v", false, "print version")
	flag.BoolVar(&debugArg, "debug", false, "debug logging")
}

func main() {
	flag.Parse()

	var loggingOpts *slog.HandlerOptions
	if debugArg {
		loggingOpts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, loggingOpts))
	slog.SetDefault(logger)

	if displayVersionArg {
		fmt.Println(version)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var ip net.IP
	port := 0
	var err error

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

	ticker := time.NewTicker(source.Refresh())
	defer ticker.Stop()
	for {
		slog.Debug("refreshing lease", "refresh", source.Refresh())
		ip, port = tick(ctx, source, actions, notifiers, ip, port)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			slog.Info("exiting")
			return
		}
	}
}

func tick(ctx context.Context, source source.Source, actions []action.Action, notifiers []notify.Notifier, oldIP net.IP, oldPort int) (net.IP, int) {
	newIP, newPort, err := source.Get()
	if err != nil {
		slog.Error("getting from source", "err", err)
		return oldIP, oldPort
	}
	if newPort != oldPort {
		if oldPort == 0 {
			slog.Info("initial port", "port", newPort)
		} else {
			slog.Info("port change", "oldPort", oldPort, "newPort", newPort)
		}
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
				slog.Error("performing action", "name", action.Name(), "err", err)
			}
		}

		slog.Info("latest endpoint", "ip", newIP, "port", newPort)
		for _, notifier := range notifiers {
			for _, result := range results {
				err = notifier.Notify(ctx, result)
				if err != nil {
					slog.Error("performing notify", "name", notifier.Name(), "err", err)
				}
			}
		}
	}

	return newIP, newPort
}

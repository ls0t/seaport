package action

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

type DuckDNS struct {
	domains string
	token   string
	txt     string
}

func NewDuckDNS(options map[string]string) (Action, error) {
	domains := options["domains"]
	if strings.TrimSpace(domains) == "" {
		return nil, errors.New("'domains' must not be empty")
	}

	token := options["token"]
	if strings.TrimSpace(token) == "" {
		token = os.Getenv("SEAPORT_DUCKDNS_TOKEN")
		if token == "" {
			return nil, errors.New("'token' must not be empty")
		}
	}

	txt := strings.TrimSpace(options["txt"])
	return &DuckDNS{
		domains: domains,
		token:   token,
		txt:     txt,
	}, nil
}

func (d *DuckDNS) Act(ctx context.Context, ip net.IP, port int) error {
	err := sendReq(ctx, ip, d.domains, d.token, "")
	if err != nil {
		return err
	}
	if d.txt != "" {
		err = sendReq(ctx, ip, d.domains, d.token, d.txt)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendReq(ctx context.Context, ip net.IP, domains string, token string, txt string) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://www.duckdns.org/update", nil)
	if err != nil {
		return fmt.Errorf("constructing http request: %w", err)
	}

	q := req.URL.Query()
	q.Add("domains", domains)
	q.Add("token", token)
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		q.Add("verbose", "true")

	}
	if txt != "" {
		q.Add("txt", txt)
	} else {
		q.Add("ip", ip.String())
	}
	req.URL.RawQuery = q.Encode()
	slog.Debug("duckdns query", "query", req.URL.RawQuery)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending http request: %w", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 100)
	n, err := resp.Body.Read(buf)
	if err != nil {
		return fmt.Errorf("reading http body: %w", err)
	}

	slog.Debug("duckdns response", "body", buf[0:n])
	if bytes.Equal(buf[0:2], []byte("KO")) {
		return fmt.Errorf("unsuccessful update: %s", buf[0:n])
	} else if !bytes.Equal(buf[0:2], []byte("OK")) {
		return fmt.Errorf("unknown response: %s", buf[0:n])
	}

	return nil
}

func (d *DuckDNS) Name() string {
	return "duckdns"
}

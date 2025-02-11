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

type FreeMyIP struct {
	domain string
	token  string
	txt    string
}

func NewFreeMyIP(options map[string]string) (Action, error) {
	domain := options["domain"]
	if strings.TrimSpace(domain) == "" {
		return nil, errors.New("'domain' must not be empty")
	}

	token := options["token"]
	if strings.TrimSpace(token) == "" {
		token = os.Getenv("SEAPORT_FREEMYIP_TOKEN")
		if token == "" {
			return nil, errors.New("'token' must not be empty")
		}
	}

	txt := options["txt"]
	return &FreeMyIP{
		domain: domain,
		token:  token,
		txt:    txt,
	}, nil
}

func (f *FreeMyIP) Act(ctx context.Context, ip net.IP, port int) error {
	err := sendFreeMyIPReq(ctx, ip, f.domain, f.token, "")
	if err != nil {
		return err
	}
	if f.txt != "" {
		err = sendFreeMyIPReq(ctx, ip, f.domain, f.token, f.txt)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendFreeMyIPReq(ctx context.Context, ip net.IP, domain string, token string, txt string) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://www.freemyip.com/update", nil)
	if err != nil {
		return fmt.Errorf("constructing http request: %w", err)
	}

	q := req.URL.Query()
	q.Add("domain", domain)
	q.Add("token", token)
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		q.Add("verbose", "true")
	}
	if txt != "" {
		q.Add("txt", txt)
	} else {
		q.Add("myip", ip.String())
	}
	req.URL.RawQuery = q.Encode()
	slog.Debug("freemyip query", "query", req.URL.RawQuery)

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

	slog.Debug("freemyip response", "body", buf[0:n])
	if bytes.Equal(buf[0:5], []byte("ERROR")) {
		return fmt.Errorf("unsuccessful update: %s", buf[0:n])
	} else if !bytes.Equal(buf[0:2], []byte("OK")) {
		return fmt.Errorf("unknown response: %s", buf[0:n])
	}

	return nil
}

func (f *FreeMyIP) Name() string {
	return "freemyip"
}

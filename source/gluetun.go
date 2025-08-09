package source

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Gluetun struct {
	url        *url.URL
	authMethod string
	username   string
	password   string
	apikey     string
}

func NewGluetun(options map[string]string) (Source, error) {
	urlStr := options["url"]
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	authMethod := strings.TrimSpace(options["auth"])
	switch authMethod {
	case "none", "":
		authMethod = ""
	case "basic", "apikey":
		authMethod = options["auth"]
	default:
		return nil, fmt.Errorf("unknown auth type: %s", authMethod)
	}

	return &Gluetun{
		url:        u,
		authMethod: authMethod,
		username:   options["username"],
		password:   options["password"],
		apikey:     options["apikey"],
	}, nil
}

func (g *Gluetun) Get() (net.IP, int, error) {
	portPath := "/v1/openvpn/portforwarded"
	resp, err := g.doReq(portPath)
	if err != nil {
		return nil, 0, fmt.Errorf("getting port from gluetun: %w", err)
	}
	port := resp.Port

	ipPath := "/v1/publicip/ip"
	resp, err = g.doReq(ipPath)
	if err != nil {
		return nil, 0, fmt.Errorf("getting IP from gluetun: %w", err)
	}
	ip := resp.PublicIP
	netIP := net.ParseIP(ip)

	return netIP, port, nil
}

func (g *Gluetun) doReq(path string) (response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", g.url.String()+path, nil)
	if err != nil {
		return response{}, fmt.Errorf("building HTTP request: %w", err)
	}

	switch g.authMethod {
	case "":
	case "basic":
		req.SetBasicAuth(g.username, g.password)
	case "apikey":
		req.Header.Set("X-API-Key", g.apikey)
	default:
		return response{}, fmt.Errorf("unknown auth method: %s", g.authMethod)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		return response{}, fmt.Errorf("decoding response failed: %w", err)
	}

	return res, nil
}

func (g *Gluetun) Refresh() time.Duration {
	return 30 * time.Second
}

type response struct {
	Port     int    `json:"port"`
	PublicIP string `json:"public_ip"`
}

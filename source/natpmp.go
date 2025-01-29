package source

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"

	natpmp "github.com/jackpal/go-nat-pmp"
)

func NewNatPMP(options map[string]string) (Source, error) {
	gatewayIP := net.ParseIP(options["gatewayIP"])
	if gatewayIP == nil {
		return nil, fmt.Errorf("gatewayIP '%s' could not be parsed", options["gatewayIP"])
	}

	var externalPort int
	var err error
	if options["externalPort"] != "" {
		externalPort, err = strconv.Atoi(options["externalPort"])
		if err != nil {
			return nil, fmt.Errorf("externalPort '%s' could not be parsed: %w", options["externalPort"], err)
		}
	}

	var internalPort int
	if options["internalPort"] != "" {
		internalPort, err = strconv.Atoi(options["internalPort"])
		if err != nil {
			return nil, fmt.Errorf("internalPort '%s' could not be parsed: %w", options["internalPort"], err)
		}
	}

	var randomPort bool
	if options["randomPort"] != "" {
		randomPort, err = strconv.ParseBool(options["randomPort"])
		if err != nil {
			return nil, fmt.Errorf("randomPort '%s' could not be parsed: %w", options["randomPort"], err)
		}
	}

	return &NatPMP{
		successfulRun: false,
		gatewayIP:     gatewayIP,
		externalPort:  externalPort,
		internalPort:  internalPort,
		randomPort:    randomPort,
	}, nil
}

type NatPMP struct {
	successfulRun bool
	gatewayIP     net.IP
	internalPort  int
	externalPort  int
	randomPort    bool
}

func (n *NatPMP) Get() (net.IP, int, error) {
	client := natpmp.NewClient(n.gatewayIP)

	requestedExternalPort := n.externalPort
	if !n.successfulRun {
		if n.randomPort {
			requestedExternalPort = rand.Intn(30000) + 30000
			n.internalPort = requestedExternalPort
		}

		if n.internalPort == 0 {
			n.internalPort = requestedExternalPort
		}
	}

	portResponse, err := client.AddPortMapping("tcp", n.internalPort, requestedExternalPort, 60)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting tcp port mapping: %w", err)
	}
	n.externalPort = int(portResponse.MappedExternalPort)

	portResponse, err = client.AddPortMapping("udp", n.internalPort, n.externalPort, 60)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting udp port mapping: %w", err)
	}
	if int(portResponse.MappedExternalPort) != n.externalPort {
		return nil, 0, fmt.Errorf("port mismatch from nat-pmp server: got %d, expected %d", portResponse.MappedExternalPort, n.externalPort)
	}

	addressResponse, err := client.GetExternalAddress()
	if err != nil {
		return nil, 0, fmt.Errorf("error getting external address: %w", err)
	}

	ip := net.IP(addressResponse.ExternalIPAddress[:])
	n.successfulRun = true

	return ip, n.externalPort, nil
}

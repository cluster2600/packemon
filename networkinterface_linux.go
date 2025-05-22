// +build linux

package packemon

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"strings"

	"golang.org/x/sys/unix"
)

// NetworkInterface represents a network interface on Linux
type NetworkInterface struct {
	Intf       *net.Interface
	Socket     int // file descriptor
	SocketAddr unix.SockaddrLinklayer
	IPAddr     uint32

	PassiveCh chan *Passive
}

// NewNetworkInterface creates a new NetworkInterface for the specified interface on Linux
func NewNetworkInterface(nwInterface string) (*NetworkInterface, error) {
	intf, err := getInterface(nwInterface)
	if err != nil {
		return nil, err
	}
	ipAddrs, err := intf.Addrs()
	if err != nil {
		return nil, err
	}

	var ipAddr uint32
	for i, addr := range ipAddrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		if i == 0 {
			ipAddr = binary.BigEndian.Uint32(ip)
		}
	}

	// Open RAW socket
	sock, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(htons(unix.ETH_P_ALL)))
	if err != nil {
		return nil, err
	}

	// Bind to interface
	addr := unix.SockaddrLinklayer{
		Protocol: htons(unix.ETH_P_ALL),
		Ifindex:  intf.Index,
	}
	// Bind socket to interface
	if err := unix.Bind(sock, &addr); err != nil {
		unix.Close(sock)
		return nil, err
	}

	nwif := &NetworkInterface{
		Intf:       intf,
		Socket:     sock,
		SocketAddr: addr,
		IPAddr:     ipAddr,
		PassiveCh:  make(chan *Passive, 100),
	}

	return nwif, nil
}

// getInterface finds the specified network interface
func getInterface(nwInterface string) (*net.Interface, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, intf := range ifs {
		if intf.Name == nwInterface || strings.Contains(intf.Name, nwInterface) {
			return &intf, nil
		}
	}

	return nil, errors.New("interface not found: " + nwInterface)
}

// SendEthernetFrame sends an Ethernet frame on Linux
func (nwif *NetworkInterface) SendEthernetFrame(ctx context.Context, data []byte) error {
	return unix.Sendto(nwif.Socket, data, 0, &nwif.SocketAddr)
}

// ReceiveEthernetFrame receives Ethernet frames on Linux
func (nwif *NetworkInterface) ReceiveEthernetFrame(ctx context.Context) {
	buf := make([]byte, 1500)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, _, err := unix.Recvfrom(nwif.Socket, buf, 0)
			if err != nil {
				continue
			}

			if n <= 14 {
				continue
			}

			frame := &EthernetFrame{
				DstAddr: buf[0:6],
				SrcAddr: buf[6:12],
				Type:    binary.BigEndian.Uint16(buf[12:14]),
				Payload: buf[14:n],
			}

			passive := &Passive{
				EthernetFrame: frame,
			}

			parseEthernetPayload(passive)

			select {
			case nwif.PassiveCh <- passive:
			default:
				// Channel is full, discard packet
			}
		}
	}
}

// GetNetworkInfo returns information about the network interface
func (nwif *NetworkInterface) GetNetworkInfo() (macAddr net.HardwareAddr, ipv4Addr net.IP, ipv6Addr net.IP) {
	ipv4 := make(net.IP, 4)
	binary.BigEndian.PutUint32(ipv4, nwif.IPAddr)
	
	return nwif.Intf.HardwareAddr, ipv4, nil
}

// Close closes the socket
func (nwif *NetworkInterface) Close() {
	if nwif.Socket != 0 {
		unix.Close(nwif.Socket)
	}
}

// htons converts a short (uint16) from host byte order to network byte order.
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

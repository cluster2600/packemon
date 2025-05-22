// +build darwin

package packemon

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// NetworkInterface represents a network interface on macOS
type NetworkInterface struct {
	Intf       *net.Interface
	Handle     *pcap.Handle
	IPAddr     uint32
	IPv6Addr   net.IP // For IPv6 support
	MacAddr    net.HardwareAddr

	PassiveCh chan *Passive
}

// NewNetworkInterface creates a new NetworkInterface for the specified interface on macOS
func NewNetworkInterface(nwInterface string) (*NetworkInterface, error) {
	intf, err := getInterface(nwInterface)
	if err != nil {
		return nil, err
	}

	// Get IP addresses associated with the interface
	ipAddrs, err := intf.Addrs()
	if err != nil {
		return nil, err
	}

	var ipAddr uint32
	var ipv6Addr net.IP

	// Find the first IPv4 and IPv6 address for the interface
	for _, addr := range ipAddrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if ip4 := ipnet.IP.To4(); ip4 != nil {
			ipAddr = binary.BigEndian.Uint32(ip4)
		} else if ipnet.IP.To16() != nil && ipAddr == 0 {
			ipv6Addr = ipnet.IP
		}
	}

	if ipAddr == 0 && ipv6Addr == nil {
		return nil, errors.New("no IP address found for interface")
	}

	// Create a new pcap handle for packet capture
	handle, err := pcap.OpenLive(intf.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("failed to open pcap handle: %v", err)
	}

	nwif := &NetworkInterface{
		Intf:      intf,
		Handle:    handle,
		IPAddr:    ipAddr,
		IPv6Addr:  ipv6Addr,
		MacAddr:   intf.HardwareAddr,
		PassiveCh: make(chan *Passive, 100),
	}

	return nwif, nil
}

// getInterface finds the specified network interface
func getInterface(nwInterface string) (*net.Interface, error) {
	// List all network interfaces
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Try to find the requested interface
	for _, intf := range ifs {
		if intf.Name == nwInterface || strings.Contains(intf.Name, nwInterface) {
			return &intf, nil
		}
	}

	return nil, fmt.Errorf("interface %s not found", nwInterface)
}

// SendEthernetFrame sends an Ethernet frame on macOS
func (nwif *NetworkInterface) SendEthernetFrame(ctx context.Context, data []byte) error {
	// Convert raw bytes to a gopacket packet to inject
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	if err := nwif.Handle.WritePacketData(data); err != nil {
		return fmt.Errorf("failed to write packet data: %v", err)
	}
	return nil
}

// ReceiveEthernetFrame receives Ethernet frames on macOS
func (nwif *NetworkInterface) ReceiveEthernetFrame(ctx context.Context) {
	packetSource := gopacket.NewPacketSource(nwif.Handle, layers.LayerTypeEthernet)
	packetChan := packetSource.Packets()

	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-packetChan:
			if packet == nil {
				continue
			}

			// Process received packet
			data := packet.Data()
			if len(data) < 14 { // Minimum Ethernet frame size
				continue
			}

			passive := &Passive{}

			// Parse Ethernet frame
			ethernetFrame := &EthernetFrame{
				DstAddr: data[0:6],
				SrcAddr: data[6:12],
				Type:    binary.BigEndian.Uint16(data[12:14]),
				Payload: data[14:],
			}
			passive.EthernetFrame = ethernetFrame

			// Parse upper-layer protocols
			parseEthernetPayload(passive)

			// Send to channel
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
	
	return nwif.MacAddr, ipv4, nwif.IPv6Addr
}

// Close cleans up resources
func (nwif *NetworkInterface) Close() {
	if nwif.Handle != nil {
		nwif.Handle.Close()
	}
}

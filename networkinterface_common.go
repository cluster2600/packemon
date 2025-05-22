package packemon

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"strings"
)

// Common interface for both Linux and macOS implementations
type NetworkInterfaceer interface {
	SendEthernetFrame(ctx context.Context, data []byte) error
	ReceiveEthernetFrame(ctx context.Context)
	GetNetworkInfo() (macAddr net.HardwareAddr, ipv4Addr net.IP, ipv6Addr net.IP)
	Close()
}

// Parse an Ethernet payload into upper-layer protocols
func parseEthernetPayload(passive *Passive) {
	if passive.EthernetFrame == nil || len(passive.EthernetFrame.Payload) == 0 {
		return
	}

	etherType := passive.EthernetFrame.Type

	switch etherType {
	case 0x0806: // ARP
		// Parse ARP packet
		if len(passive.EthernetFrame.Payload) >= 28 {
			// Minimum ARP packet size
			arp := ParseARPPacket(passive.EthernetFrame.Payload)
			passive.ARP = arp
		}

	case 0x0800: // IPv4
		// Parse IPv4 packet
		if len(passive.EthernetFrame.Payload) >= 20 {
			// Minimum IPv4 header size
			ipv4 := ParseIPv4Packet(passive.EthernetFrame.Payload)
			passive.IPv4 = ipv4

			// Parse upper layer based on protocol
			if ipv4 != nil && len(ipv4.Payload) > 0 {
				parseIPv4Payload(passive, ipv4)
			}
		}

	case 0x86DD: // IPv6
		// Parse IPv6 packet
		if len(passive.EthernetFrame.Payload) >= 40 {
			// IPv6 header size
			ipv6 := ParseIPv6Packet(passive.EthernetFrame.Payload)
			passive.IPv6 = ipv6

			// Parse upper layer based on next header
			if ipv6 != nil && len(ipv6.Payload) > 0 {
				parseIPv6Payload(passive, ipv6)
			}
		}
	}
}

// Parse an IPv4 payload into upper-layer protocols
func parseIPv4Payload(passive *Passive, ipv4 *IPv4Packet) {
	switch ipv4.Protocol {
	case 1: // ICMP
		if len(ipv4.Payload) >= 8 {
			// Minimum ICMP message size
			icmp := ParseICMPPacket(ipv4.Payload)
			passive.ICMP = icmp
		}

	case 6: // TCP
		if len(ipv4.Payload) >= 20 {
			// Minimum TCP header size
			tcp := ParseTCPPacket(ipv4.Payload)
			passive.TCP = tcp

			// Parse application layer protocols based on port
			if tcp != nil && len(tcp.Payload) > 0 {
				parseTCPPayload(passive, tcp)
			}
		}

	case 17: // UDP
		if len(ipv4.Payload) >= 8 {
			// UDP header size
			udp := ParseUDPPacket(ipv4.Payload)
			passive.UDP = udp

			// Parse application layer protocols based on port
			if udp != nil && len(udp.Payload) > 0 {
				parseUDPPayload(passive, udp)
			}
		}
	}
}

// Parse an IPv6 payload into upper-layer protocols
func parseIPv6Payload(passive *Passive, ipv6 *IPv6Packet) {
	switch ipv6.NextHeader {
	case 58: // ICMPv6
		if len(ipv6.Payload) >= 8 {
			// Minimum ICMPv6 message size
			icmpv6 := ParseICMPv6Packet(ipv6.Payload)
			passive.ICMPv6 = icmpv6
		}

	case 6: // TCP
		if len(ipv6.Payload) >= 20 {
			// Minimum TCP header size
			tcp := ParseTCPPacket(ipv6.Payload)
			passive.TCP = tcp

			// Parse application layer protocols based on port
			if tcp != nil && len(tcp.Payload) > 0 {
				parseTCPPayload(passive, tcp)
			}
		}

	case 17: // UDP
		if len(ipv6.Payload) >= 8 {
			// UDP header size
			udp := ParseUDPPacket(ipv6.Payload)
			passive.UDP = udp

			// Parse application layer protocols based on port
			if udp != nil && len(udp.Payload) > 0 {
				parseUDPPayload(passive, udp)
			}
		}
	}
}

// Parse TCP payload based on port numbers
func parseTCPPayload(passive *Passive, tcp *TCPPacket) {
	// HTTP (port 80)
	if tcp.DstPort == 80 || tcp.SrcPort == 80 {
		if tcp.DstPort == 80 {
			// HTTP Request
			http := ParseHTTPRequest(tcp.Payload)
			if http != nil {
				passive.HTTP = http
			}
		} else {
			// HTTP Response
			httpRes := ParseHTTPResponse(tcp.Payload)
			if httpRes != nil {
				passive.HTTPRes = httpRes
			}
		}
	}

	// HTTPS (port 443)
	if tcp.DstPort == 443 || tcp.SrcPort == 443 {
		// TLS parsing
		ParseTLSData(tcp.Payload, passive)
	}

	// DNS over TCP (port 53)
	if tcp.DstPort == 53 || tcp.SrcPort == 53 {
		if len(tcp.Payload) > 2 {
			// Skip TCP DNS length field (first 2 bytes)
			dnsData := tcp.Payload[2:]
			parseDNSData(dnsData, passive)
		}
	}
}

// Parse UDP payload based on port numbers
func parseUDPPayload(passive *Passive, udp *UDPPacket) {
	// DNS (port 53)
	if udp.DstPort == 53 || udp.SrcPort == 53 {
		parseDNSData(udp.Payload, passive)
	}
}

// Parse DNS data
func parseDNSData(data []byte, passive *Passive) {
	if len(data) < 12 {
		// DNS header is 12 bytes
		return
	}

	// Check if it's a DNS query or response
	flags := binary.BigEndian.Uint16(data[2:4])
	isResponse := (flags & 0x8000) != 0

	if isResponse {
		dns := ParseDNSResponse(data)
		passive.DNS = dns
	} else {
		dns := ParseDNSRequest(data)
		passive.DNS = dns
	}
}

// Function to determine if the interface name is valid
func isValidInterfaceName(name string) bool {
	if name == "" {
		return false
	}
	
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	
	for _, iface := range interfaces {
		if iface.Name == name || strings.Contains(iface.Name, name) {
			return true
		}
	}
	
	return false
}

// htons converts a short (uint16) from host byte order to network byte order.
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

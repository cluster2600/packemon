// passive.go contains structures for parsed packets

package packemon

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Passive represents a parsed packet with all layers
type Passive struct {
	EthernetFrame *EthernetFrame
	ARP           *ARPPacket
	IPv4          *IPv4Packet
	IPv6          *IPv6Packet
	ICMP          *ICMPPacket
	ICMPv6        *ICMPv6Packet
	TCP           *TCPPacket
	UDP           *UDPPacket
	TLS           *TLSRecord
	DNS           *DNSPacket
	HTTP          *HTTPRequest
	HTTPRes       *HTTPResponse
}

// EthernetFrame represents an Ethernet frame
type EthernetFrame struct {
	DstAddr []byte
	SrcAddr []byte
	Type    uint16
	Payload []byte
}

// String returns a string representation of the Ethernet frame
func (e *EthernetFrame) String() string {
	return fmt.Sprintf("Ethernet Frame: Dst=%s, Src=%s, Type=0x%04x, Len=%d",
		net.HardwareAddr(e.DstAddr),
		net.HardwareAddr(e.SrcAddr),
		e.Type,
		len(e.Payload))
}

// ARPPacket represents an ARP packet
type ARPPacket struct {
	HardwareType    uint16
	ProtocolType    uint16
	HardwareSize    uint8
	ProtocolSize    uint8
	Operation       uint16
	SenderMAC       []byte
	SenderIP        []byte
	TargetMAC       []byte
	TargetIP        []byte
}

// String returns a string representation of the ARP packet
func (a *ARPPacket) String() string {
	return fmt.Sprintf("ARP: Op=%d, Sender=%s/%s, Target=%s/%s",
		a.Operation,
		net.HardwareAddr(a.SenderMAC),
		net.IP(a.SenderIP),
		net.HardwareAddr(a.TargetMAC),
		net.IP(a.TargetIP))
}

// IPv4Packet represents an IPv4 packet
type IPv4Packet struct {
	Version     uint8
	IHL         uint8
	TOS         uint8
	TotalLength uint16
	ID          uint16
	Flags       uint8
	FragOffset  uint16
	TTL         uint8
	Protocol    uint8
	Checksum    uint16
	SrcIP       []byte
	DstIP       []byte
	Options     []byte
	Payload     []byte
}

// String returns a string representation of the IPv4 packet
func (i *IPv4Packet) String() string {
	return fmt.Sprintf("IPv4: Src=%s, Dst=%s, Proto=%d, Len=%d",
		net.IP(i.SrcIP),
		net.IP(i.DstIP),
		i.Protocol,
		len(i.Payload))
}

// IPv6Packet represents an IPv6 packet
type IPv6Packet struct {
	Version      uint8
	TrafficClass uint8
	FlowLabel    uint32
	PayloadLen   uint16
	NextHeader   uint8
	HopLimit     uint8
	SrcIP        []byte
	DstIP        []byte
	Payload      []byte
}

// String returns a string representation of the IPv6 packet
func (i *IPv6Packet) String() string {
	return fmt.Sprintf("IPv6: Src=%s, Dst=%s, NextHeader=%d, Len=%d",
		net.IP(i.SrcIP),
		net.IP(i.DstIP),
		i.NextHeader,
		len(i.Payload))
}

// ICMPPacket represents an ICMP packet
type ICMPPacket struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Sequence uint16
	Payload  []byte
}

// String returns a string representation of the ICMP packet
func (i *ICMPPacket) String() string {
	return fmt.Sprintf("ICMP: Type=%d, Code=%d, ID=%d, Seq=%d",
		i.Type,
		i.Code,
		i.ID,
		i.Sequence)
}

// ICMPv6Packet represents an ICMPv6 packet
type ICMPv6Packet struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Payload  []byte
}

// String returns a string representation of the ICMPv6 packet
func (i *ICMPv6Packet) String() string {
	return fmt.Sprintf("ICMPv6: Type=%d, Code=%d",
		i.Type,
		i.Code)
}

// TCPPacket represents a TCP packet
type TCPPacket struct {
	SrcPort    uint16
	DstPort    uint16
	SeqNum     uint32
	AckNum     uint32
	DataOffset uint8
	Flags      uint8
	Window     uint16
	Checksum   uint16
	UrgPtr     uint16
	Options    []byte
	Payload    []byte
}

// String returns a string representation of the TCP packet
func (t *TCPPacket) String() string {
	return fmt.Sprintf("TCP: Src=%d, Dst=%d, Seq=%d, Ack=%d, Flags=0x%02x",
		t.SrcPort,
		t.DstPort,
		t.SeqNum,
		t.AckNum,
		t.Flags)
}

// UDPPacket represents a UDP packet
type UDPPacket struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
	Payload  []byte
}

// String returns a string representation of the UDP packet
func (u *UDPPacket) String() string {
	return fmt.Sprintf("UDP: Src=%d, Dst=%d, Len=%d",
		u.SrcPort,
		u.DstPort,
		u.Length)
}

// TLSRecord represents a TLS record
type TLSRecord struct {
	Type    uint8
	Version uint16
	Length  uint16
	Data    []byte
}

// String returns a string representation of the TLS record
func (t *TLSRecord) String() string {
	return fmt.Sprintf("TLS: Type=%d, Version=0x%04x, Len=%d",
		t.Type,
		t.Version,
		t.Length)
}

// DNSPacket represents a DNS packet
type DNSPacket struct {
	ID            uint16
	Flags         uint16
	Questions     uint16
	AnswerRRs     uint16
	AuthorityRRs  uint16
	AdditionalRRs uint16
	Payload       []byte
}

// String returns a string representation of the DNS packet
func (d *DNSPacket) String() string {
	return fmt.Sprintf("DNS: ID=%d, Flags=0x%04x, Questions=%d, Answers=%d",
		d.ID,
		d.Flags,
		d.Questions,
		d.AnswerRRs)
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string
	URI     string
	Version string
	Headers map[string]string
	Body    []byte
}

// String returns a string representation of the HTTP request
func (h *HTTPRequest) String() string {
	return fmt.Sprintf("HTTP Request: %s %s %s",
		h.Method,
		h.URI,
		h.Version)
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	Version    string
	StatusCode int
	Status     string
	Headers    map[string]string
	Body       []byte
}

// String returns a string representation of the HTTP response
func (h *HTTPResponse) String() string {
	return fmt.Sprintf("HTTP Response: %s %d %s",
		h.Version,
		h.StatusCode,
		h.Status)
}

// Parse functions

// ParseARPPacket parses ARP packet data
func ParseARPPacket(data []byte) *ARPPacket {
	if len(data) < 28 {
		return nil
	}
	
	return &ARPPacket{
		HardwareType: binary.BigEndian.Uint16(data[0:2]),
		ProtocolType: binary.BigEndian.Uint16(data[2:4]),
		HardwareSize: data[4],
		ProtocolSize: data[5],
		Operation:    binary.BigEndian.Uint16(data[6:8]),
		SenderMAC:    data[8:14],
		SenderIP:     data[14:18],
		TargetMAC:    data[18:24],
		TargetIP:     data[24:28],
	}
}

// ParseIPv4Packet parses IPv4 packet data
func ParseIPv4Packet(data []byte) *IPv4Packet {
	if len(data) < 20 {
		return nil
	}
	
	ihl := (data[0] & 0x0F) * 4
	if len(data) < int(ihl) {
		return nil
	}
	
	return &IPv4Packet{
		Version:     (data[0] >> 4) & 0x0F,
		IHL:         ihl,
		TOS:         data[1],
		TotalLength: binary.BigEndian.Uint16(data[2:4]),
		ID:          binary.BigEndian.Uint16(data[4:6]),
		Flags:       (data[6] >> 5) & 0x07,
		FragOffset:  binary.BigEndian.Uint16(data[6:8]) & 0x1FFF,
		TTL:         data[8],
		Protocol:    data[9],
		Checksum:    binary.BigEndian.Uint16(data[10:12]),
		SrcIP:       data[12:16],
		DstIP:       data[16:20],
		Options:     data[20:ihl],
		Payload:     data[ihl:],
	}
}

// ParseIPv6Packet parses IPv6 packet data
func ParseIPv6Packet(data []byte) *IPv6Packet {
	if len(data) < 40 {
		return nil
	}
	
	return &IPv6Packet{
		Version:      (data[0] >> 4) & 0x0F,
		TrafficClass: ((data[0] & 0x0F) << 4) | ((data[1] >> 4) & 0x0F),
		FlowLabel:    uint32(data[1]&0x0F)<<16 | uint32(data[2])<<8 | uint32(data[3]),
		PayloadLen:   binary.BigEndian.Uint16(data[4:6]),
		NextHeader:   data[6],
		HopLimit:     data[7],
		SrcIP:        data[8:24],
		DstIP:        data[24:40],
		Payload:      data[40:],
	}
}

// ParseICMPPacket parses ICMP packet data
func ParseICMPPacket(data []byte) *ICMPPacket {
	if len(data) < 8 {
		return nil
	}
	
	return &ICMPPacket{
		Type:     data[0],
		Code:     data[1],
		Checksum: binary.BigEndian.Uint16(data[2:4]),
		ID:       binary.BigEndian.Uint16(data[4:6]),
		Sequence: binary.BigEndian.Uint16(data[6:8]),
		Payload:  data[8:],
	}
}

// ParseICMPv6Packet parses ICMPv6 packet data
func ParseICMPv6Packet(data []byte) *ICMPv6Packet {
	if len(data) < 4 {
		return nil
	}
	
	return &ICMPv6Packet{
		Type:     data[0],
		Code:     data[1],
		Checksum: binary.BigEndian.Uint16(data[2:4]),
		Payload:  data[4:],
	}
}

// ParseTCPPacket parses TCP packet data
func ParseTCPPacket(data []byte) *TCPPacket {
	if len(data) < 20 {
		return nil
	}
	
	dataOffset := (data[12] >> 4) * 4
	if len(data) < int(dataOffset) {
		return nil
	}
	
	return &TCPPacket{
		SrcPort:    binary.BigEndian.Uint16(data[0:2]),
		DstPort:    binary.BigEndian.Uint16(data[2:4]),
		SeqNum:     binary.BigEndian.Uint32(data[4:8]),
		AckNum:     binary.BigEndian.Uint32(data[8:12]),
		DataOffset: dataOffset,
		Flags:      data[13],
		Window:     binary.BigEndian.Uint16(data[14:16]),
		Checksum:   binary.BigEndian.Uint16(data[16:18]),
		UrgPtr:     binary.BigEndian.Uint16(data[18:20]),
		Options:    data[20:dataOffset],
		Payload:    data[dataOffset:],
	}
}

// ParseUDPPacket parses UDP packet data
func ParseUDPPacket(data []byte) *UDPPacket {
	if len(data) < 8 {
		return nil
	}
	
	return &UDPPacket{
		SrcPort:  binary.BigEndian.Uint16(data[0:2]),
		DstPort:  binary.BigEndian.Uint16(data[2:4]),
		Length:   binary.BigEndian.Uint16(data[4:6]),
		Checksum: binary.BigEndian.Uint16(data[6:8]),
		Payload:  data[8:],
	}
}

// ParseDNSRequest parses DNS request data
func ParseDNSRequest(data []byte) *DNSPacket {
	if len(data) < 12 {
		return nil
	}
	
	return &DNSPacket{
		ID:            binary.BigEndian.Uint16(data[0:2]),
		Flags:         binary.BigEndian.Uint16(data[2:4]),
		Questions:     binary.BigEndian.Uint16(data[4:6]),
		AnswerRRs:     binary.BigEndian.Uint16(data[6:8]),
		AuthorityRRs:  binary.BigEndian.Uint16(data[8:10]),
		AdditionalRRs: binary.BigEndian.Uint16(data[10:12]),
		Payload:       data[12:],
	}
}

// ParseDNSResponse parses DNS response data
func ParseDNSResponse(data []byte) *DNSPacket {
	return ParseDNSRequest(data) // Same structure
}

// ParseHTTPRequest parses HTTP request data
func ParseHTTPRequest(data []byte) *HTTPRequest {
	// Simplified implementation
	return &HTTPRequest{
		Method:  "GET", // Placeholder
		URI:     "/",   // Placeholder
		Version: "HTTP/1.1",
		Headers: make(map[string]string),
		Body:    []byte{},
	}
}

// ParseHTTPResponse parses HTTP response data
func ParseHTTPResponse(data []byte) *HTTPResponse {
	// Simplified implementation
	return &HTTPResponse{
		Version:    "HTTP/1.1",
		StatusCode: 200,
		Status:     "OK",
		Headers:    make(map[string]string),
		Body:       []byte{},
	}
}

// ParseTLSData parses TLS data
func ParseTLSData(data []byte, passive *Passive) {
	if len(data) < 5 {
		return
	}
	
	passive.TLS = &TLSRecord{
		Type:    data[0],
		Version: binary.BigEndian.Uint16(data[1:3]),
		Length:  binary.BigEndian.Uint16(data[3:5]),
		Data:    data[5:],
	}
}

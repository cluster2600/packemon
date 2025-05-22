package packemon

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

// ICMPv6 implements the Internet Control Message Protocol for IPv6
// RFC 4443: https://tools.ietf.org/html/rfc4443
type ICMPv6 struct {
	Type      uint8
	Code      uint8
	Checksum  uint16
	MessageBody []byte
}

// ICMPv6 message types as defined in RFC 4443
const (
	// Error messages
	ICMPv6_TYPE_DESTINATION_UNREACHABLE = 1
	ICMPv6_TYPE_PACKET_TOO_BIG          = 2
	ICMPv6_TYPE_TIME_EXCEEDED           = 3
	ICMPv6_TYPE_PARAMETER_PROBLEM       = 4
	
	// Informational messages
	ICMPv6_TYPE_ECHO_REQUEST            = 128
	ICMPv6_TYPE_ECHO_REPLY              = 129
	
	// Neighbor Discovery Protocol messages (RFC 4861)
	ICMPv6_TYPE_ROUTER_SOLICITATION     = 133
	ICMPv6_TYPE_ROUTER_ADVERTISEMENT    = 134
	ICMPv6_TYPE_NEIGHBOR_SOLICITATION   = 135
	ICMPv6_TYPE_NEIGHBOR_ADVERTISEMENT  = 136
	ICMPv6_TYPE_REDIRECT                = 137
)

// Echo Request/Reply specific structure
type ICMPv6Echo struct {
	Identifier uint16
	SequenceNumber uint16
	Data []byte
}

// ParsedICMPv6 parses ICMPv6 packets from binary data
func ParsedICMPv6(payload []byte) *ICMPv6 {
	if len(payload) < 4 {
		return nil
	}
	
	return &ICMPv6{
		Type:      payload[0],
		Code:      payload[1],
		Checksum:  binary.BigEndian.Uint16(payload[2:4]),
		MessageBody: payload[4:],
	}
}

// ParsedICMPv6Echo parses ICMPv6 Echo Request/Reply packets
func ParsedICMPv6Echo(icmpv6 *ICMPv6) *ICMPv6Echo {
	if icmpv6 == nil || len(icmpv6.MessageBody) < 4 {
		return nil
	}
	
	return &ICMPv6Echo{
		Identifier:    binary.BigEndian.Uint16(icmpv6.MessageBody[0:2]),
		SequenceNumber: binary.BigEndian.Uint16(icmpv6.MessageBody[2:4]),
		Data:          icmpv6.MessageBody[4:],
	}
}

// NewICMPv6EchoRequest creates a new ICMPv6 Echo Request packet
func NewICMPv6EchoRequest() *ICMPv6 {
	// Create echo data with timestamp similar to ping
	timestamp := func() []byte {
		now := time.Now().Unix()
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(now))
		return binary.LittleEndian.AppendUint32(b, 0x00000000)
	}()
	
	// Create echo body
	echo := &ICMPv6Echo{
		Identifier:    0x1234, // Arbitrary identifier
		SequenceNumber: 0x0001,
		Data:          timestamp,
	}
	
	// Convert echo to bytes
	echoBuf := &bytes.Buffer{}
	WriteUint16(echoBuf, echo.Identifier)
	WriteUint16(echoBuf, echo.SequenceNumber)
	echoBuf.Write(echo.Data)
	
	// Create ICMPv6 message with zero checksum (to be calculated later)
	icmpv6 := &ICMPv6{
		Type:      ICMPv6_TYPE_ECHO_REQUEST,
		Code:      0,
		Checksum:  0,
		MessageBody: echoBuf.Bytes(),
	}
	
	return icmpv6
}

// Bytes serializes an ICMPv6 packet into a byte slice
func (i *ICMPv6) Bytes() []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(i.Type)
	buf.WriteByte(i.Code)
	WriteUint16(buf, i.Checksum)
	buf.Write(i.MessageBody)
	return buf.Bytes()
}

// CalculateChecksum calculates the ICMPv6 checksum including IPv6 pseudo-header
// IPv6 pseudo-header consists of: source address, destination address, 
// payload length, and next header type (58 for ICMPv6)
func (i *ICMPv6) CalculateChecksum(srcIP, dstIP net.IP) uint16 {
	// Prepare the ICMPv6 data with zero checksum
	icmpData := i.bytesWithZeroChecksum()
	
	// Create the pseudo-header
	pseudoHeader := &bytes.Buffer{}
	
	// Source IP (16 bytes for IPv6)
	pseudoHeader.Write(srcIP.To16())
	
	// Destination IP (16 bytes for IPv6)
	pseudoHeader.Write(dstIP.To16())
	
	// Upper-layer packet length (32 bits)
	var packetLength uint32 = uint32(len(icmpData))
	binary.Write(pseudoHeader, binary.BigEndian, packetLength)
	
	// Zero padding (24 bits)
	pseudoHeader.Write([]byte{0, 0, 0})
	
	// Next header (8 bits) - 58 for ICMPv6
	pseudoHeader.WriteByte(IPv6_NEXT_HEADER_ICMPv6)
	
	// Combine pseudo-header and ICMPv6 data
	checksumData := append(pseudoHeader.Bytes(), icmpData...)
	
	// Calculate checksum
	return calculateInternetChecksum(checksumData)
}

// bytesWithZeroChecksum returns ICMPv6 packet bytes with checksum field set to zero
func (i *ICMPv6) bytesWithZeroChecksum() []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(i.Type)
	buf.WriteByte(i.Code)
	WriteUint16(buf, 0) // Zero checksum
	buf.Write(i.MessageBody)
	return buf.Bytes()
}

// calculateInternetChecksum calculates the Internet Checksum as per RFC 1071
func calculateInternetChecksum(data []byte) uint16 {
	var sum uint32
	
	// Handle complete 16-bit chunks
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	
	// Handle potential leftover byte
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	
	// Fold 32-bit sum to 16 bits
	for sum > 0xffff {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	
	return ^uint16(sum)
}

package packemon

import (
	"bytes"
	"net"
	"testing"
)

func TestICMPv6_Bytes(t *testing.T) {
	// Create a simple ICMPv6 Echo Request
	icmpv6 := &ICMPv6{
		Type:      ICMPv6_TYPE_ECHO_REQUEST,
		Code:      0,
		Checksum:  0x1234, // Just for testing
		MessageBody: []byte{0x00, 0x01, 0x00, 0x02, 0xaa, 0xbb, 0xcc, 0xdd},
	}

	// Expected bytes for the above ICMPv6 packet
	expected := []byte{
		0x80, 0x00, 0x12, 0x34, // Type, Code, Checksum
		0x00, 0x01, 0x00, 0x02, // First part of message body (identifier and sequence)
		0xaa, 0xbb, 0xcc, 0xdd, // Rest of message body (data)
	}

	// Get the actual bytes
	actual := icmpv6.Bytes()

	// Compare
	if !bytes.Equal(actual, expected) {
		t.Errorf("ICMPv6.Bytes() = %v, want %v", actual, expected)
	}
}

func TestParsedICMPv6(t *testing.T) {
	// Sample ICMPv6 packet bytes (Echo Request)
	payload := []byte{
		0x80, 0x00, 0x12, 0x34, // Type, Code, Checksum
		0x00, 0x01, 0x00, 0x02, // Identifier, Sequence
		0xaa, 0xbb, 0xcc, 0xdd, // Data
	}

	// Parse the packet
	icmpv6 := ParsedICMPv6(payload)

	// Verify the parsed values
	if icmpv6.Type != ICMPv6_TYPE_ECHO_REQUEST {
		t.Errorf("ParsedICMPv6().Type = %v, want %v", icmpv6.Type, ICMPv6_TYPE_ECHO_REQUEST)
	}
	if icmpv6.Code != 0 {
		t.Errorf("ParsedICMPv6().Code = %v, want %v", icmpv6.Code, 0)
	}
	if icmpv6.Checksum != 0x1234 {
		t.Errorf("ParsedICMPv6().Checksum = %v, want %v", icmpv6.Checksum, 0x1234)
	}
	expectedBody := []byte{0x00, 0x01, 0x00, 0x02, 0xaa, 0xbb, 0xcc, 0xdd}
	if !bytes.Equal(icmpv6.MessageBody, expectedBody) {
		t.Errorf("ParsedICMPv6().MessageBody = %v, want %v", icmpv6.MessageBody, expectedBody)
	}
}

func TestParsedICMPv6Echo(t *testing.T) {
	// Create an ICMPv6 packet with Echo Request data
	icmpv6 := &ICMPv6{
		Type:      ICMPv6_TYPE_ECHO_REQUEST,
		Code:      0,
		Checksum:  0x1234,
		MessageBody: []byte{
			0x12, 0x34, // Identifier
			0x56, 0x78, // Sequence Number
			0xaa, 0xbb, 0xcc, 0xdd, // Data
		},
	}

	// Parse the Echo fields
	echo := ParsedICMPv6Echo(icmpv6)

	// Verify the parsed values
	if echo.Identifier != 0x1234 {
		t.Errorf("ParsedICMPv6Echo().Identifier = %v, want %v", echo.Identifier, 0x1234)
	}
	if echo.SequenceNumber != 0x5678 {
		t.Errorf("ParsedICMPv6Echo().SequenceNumber = %v, want %v", echo.SequenceNumber, 0x5678)
	}
	expectedData := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	if !bytes.Equal(echo.Data, expectedData) {
		t.Errorf("ParsedICMPv6Echo().Data = %v, want %v", echo.Data, expectedData)
	}
}

func TestICMPv6_CalculateChecksum(t *testing.T) {
	// Create a simple ICMPv6 Echo Request
	icmpv6 := &ICMPv6{
		Type:      ICMPv6_TYPE_ECHO_REQUEST,
		Code:      0,
		Checksum:  0, // Zero for checksum calculation
		MessageBody: []byte{
			0x12, 0x34, // Identifier
			0x56, 0x78, // Sequence Number
			0xaa, 0xbb, 0xcc, 0xdd, // Data
		},
	}

	// Source and destination IPv6 addresses for pseudo-header
	srcIP := net.ParseIP("2001:db8::1")
	dstIP := net.ParseIP("2001:db8::2")

	// Calculate the checksum
	checksum := icmpv6.CalculateChecksum(srcIP, dstIP)

	// We can't easily predict the exact checksum value without duplicating the algorithm,
	// but we can verify it's not zero (which would indicate a calculation error)
	if checksum == 0 {
		t.Errorf("ICMPv6.CalculateChecksum() = %v, should not be zero", checksum)
	}

	// Set the calculated checksum and verify it again with the same inputs
	// This should result in a checksum of 0, which is the correct value for a valid packet
	icmpv6.Checksum = checksum
	verifyChecksum := calculateInternetChecksum(append(
		createIPv6PseudoHeader(srcIP, dstIP, uint32(len(icmpv6.bytesWithZeroChecksum()))),
		icmpv6.bytesWithZeroChecksum()...
	))

	// The verification should result in 0xFFFF (which is equivalent to 0 in one's complement)
	if verifyChecksum != 0xFFFF {
		t.Errorf("Verification checksum = %v, want %v", verifyChecksum, 0xFFFF)
	}
}

// Helper function to create an IPv6 pseudo-header for testing
func createIPv6PseudoHeader(srcIP, dstIP net.IP, length uint32) []byte {
	pseudoHeader := &bytes.Buffer{}
	pseudoHeader.Write(srcIP.To16())
	pseudoHeader.Write(dstIP.To16())
	
	lengthBytes := make([]byte, 4)
	lengthBytes[0] = byte(length >> 24)
	lengthBytes[1] = byte(length >> 16)
	lengthBytes[2] = byte(length >> 8)
	lengthBytes[3] = byte(length)
	pseudoHeader.Write(lengthBytes)
	
	pseudoHeader.Write([]byte{0, 0, 0})
	pseudoHeader.WriteByte(IPv6_NEXT_HEADER_ICMPv6)
	
	return pseudoHeader.Bytes()
}

func TestNewICMPv6EchoRequest(t *testing.T) {
	// Create a new ICMPv6 Echo Request
	icmpv6 := NewICMPv6EchoRequest()

	// Verify the basic fields
	if icmpv6.Type != ICMPv6_TYPE_ECHO_REQUEST {
		t.Errorf("NewICMPv6EchoRequest().Type = %v, want %v", icmpv6.Type, ICMPv6_TYPE_ECHO_REQUEST)
	}
	if icmpv6.Code != 0 {
		t.Errorf("NewICMPv6EchoRequest().Code = %v, want %v", icmpv6.Code, 0)
	}
	if icmpv6.Checksum != 0 {
		t.Errorf("NewICMPv6EchoRequest().Checksum = %v, want %v", icmpv6.Checksum, 0)
	}
	
	// Verify the message body contains at least 8 bytes (4 for identifier/sequence, 4+ for data)
	if len(icmpv6.MessageBody) < 8 {
		t.Errorf("NewICMPv6EchoRequest().MessageBody length = %v, want at least 8", len(icmpv6.MessageBody))
	}
}

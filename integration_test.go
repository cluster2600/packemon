package packemon

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// Integration tests for packet generation and monitoring
// These tests verify that packets can be properly generated, sent, and monitored

// TestPacketGenerationAndMonitoring tests the full cycle of generating a packet,
// sending it, and monitoring it on the network interface
func TestPacketGenerationAndMonitoring(t *testing.T) {
	// Skip this test in short mode as it requires network access
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a network interface for testing
	// Use the loopback interface to avoid sending packets to external networks
	intf, err := NewNetworkInterface("lo")
	if err != nil {
		t.Fatalf("Failed to create network interface: %v", err)
	}
	defer intf.Close()

	// Start monitoring in a goroutine
	monitorDone := make(chan struct{})
	packetReceived := make(chan *Passive)
	
	go func() {
		defer close(monitorDone)
		
		// Start receiving packets
		go func() {
			for passive := range intf.PassiveCh {
				// Only process our test packets (ICMPv6 Echo Request)
				if passive.ICMPv6 != nil && passive.ICMPv6.Type == ICMPv6_TYPE_ECHO_REQUEST {
					packetReceived <- passive
					return
				}
			}
		}()
		
		// Set up a timeout
		timeout := time.After(5 * time.Second)
		select {
		case <-timeout:
			t.Error("Timeout waiting for packet")
		case <-packetReceived:
			// Successfully received the packet
		}
	}()

	// Create an ICMPv6 Echo Request packet
	icmpv6 := NewICMPv6EchoRequest()
	
	// Create an IPv6 packet with the ICMPv6 payload
	ipv6 := &IPv6{
		Version:      0x06,
		TrafficClass: 0x00,
		FlowLabel:    0x00000,
		PayloadLength: uint16(len(icmpv6.Bytes())),
		NextHeader:   IPv6_NEXT_HEADER_ICMPv6,
		HopLimit:     64,
		SrcAddr:      net.ParseIP("::1").To16(), // Loopback source
		DstAddr:      net.ParseIP("::1").To16(), // Loopback destination
		Data:         icmpv6.Bytes(),
	}
	
	// Calculate ICMPv6 checksum
	icmpv6.Checksum = icmpv6.CalculateChecksum(ipv6.SrcAddr, ipv6.DstAddr)
	ipv6.Data = icmpv6.Bytes() // Update with correct checksum
	
	// Create an Ethernet frame with the IPv6 payload
	ethernetFrame := &EthernetFrame{
		Header: &EthernetHeader{
			Dst: HardwareAddr(intf.Intf.HardwareAddr), // Send to our own MAC
			Src: HardwareAddr(intf.Intf.HardwareAddr), // From our own MAC
			Typ: ETHER_TYPE_IPv6,
		},
		Data: ipv6.Bytes(),
	}
	
	// Send the packet
	err = intf.Send(ethernetFrame)
	if err != nil {
		t.Fatalf("Failed to send packet: %v", err)
	}
	
	// Wait for monitoring to complete
	<-monitorDone
}

// TestICMPv6EchoRequestResponse tests sending an ICMPv6 Echo Request and receiving a response
func TestICMPv6EchoRequestResponse(t *testing.T) {
	// Skip this test in short mode as it requires network access
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a network interface for testing
	intf, err := NewNetworkInterface("lo")
	if err != nil {
		t.Fatalf("Failed to create network interface: %v", err)
	}
	defer intf.Close()

	// Start monitoring in a goroutine
	monitorDone := make(chan struct{})
	echoReply := make(chan *ICMPv6)
	
	go func() {
		defer close(monitorDone)
		
		// Start receiving packets
		go func() {
			for passive := range intf.PassiveCh {
				// Look for ICMPv6 Echo Reply packets
				if passive.ICMPv6 != nil && passive.ICMPv6.Type == ICMPv6_TYPE_ECHO_REPLY {
					echoReply <- passive.ICMPv6
					return
				}
			}
		}()
		
		// Set up a timeout
		timeout := time.After(5 * time.Second)
		select {
		case <-timeout:
			t.Error("Timeout waiting for Echo Reply")
		case <-echoReply:
			// Successfully received the Echo Reply
		}
	}()

	// Create a unique identifier and sequence number for our Echo Request
	identifier := uint16(0xABCD)
	sequenceNumber := uint16(0x1234)
	
	// Create the Echo Request message body
	echoData := []byte{0xDE, 0xAD, 0xBE, 0xEF} // Unique data pattern
	echoBuf := &bytes.Buffer{}
	WriteUint16(echoBuf, identifier)
	WriteUint16(echoBuf, sequenceNumber)
	echoBuf.Write(echoData)
	
	// Create an ICMPv6 Echo Request packet
	icmpv6 := &ICMPv6{
		Type:      ICMPv6_TYPE_ECHO_REQUEST,
		Code:      0,
		Checksum:  0, // Will be calculated later
		MessageBody: echoBuf.Bytes(),
	}
	
	// Create an IPv6 packet with the ICMPv6 payload
	ipv6 := &IPv6{
		Version:      0x06,
		TrafficClass: 0x00,
		FlowLabel:    0x00000,
		PayloadLength: uint16(len(icmpv6.Bytes())),
		NextHeader:   IPv6_NEXT_HEADER_ICMPv6,
		HopLimit:     64,
		SrcAddr:      net.ParseIP("::1").To16(), // Loopback source
		DstAddr:      net.ParseIP("::1").To16(), // Loopback destination
		Data:         icmpv6.Bytes(),
	}
	
	// Calculate ICMPv6 checksum
	icmpv6.Checksum = icmpv6.CalculateChecksum(ipv6.SrcAddr, ipv6.DstAddr)
	ipv6.Data = icmpv6.Bytes() // Update with correct checksum
	
	// Create an Ethernet frame with the IPv6 payload
	ethernetFrame := &EthernetFrame{
		Header: &EthernetHeader{
			Dst: HardwareAddr(intf.Intf.HardwareAddr), // Send to our own MAC
			Src: HardwareAddr(intf.Intf.HardwareAddr), // From our own MAC
			Typ: ETHER_TYPE_IPv6,
		},
		Data: ipv6.Bytes(),
	}
	
	// Send the packet
	err = intf.Send(ethernetFrame)
	if err != nil {
		t.Fatalf("Failed to send packet: %v", err)
	}
	
	// Wait for monitoring to complete
	<-monitorDone
	
	// Verify the Echo Reply (if we got one)
	reply := <-echoReply
	if reply == nil {
		t.Fatal("Did not receive Echo Reply")
	}
	
	// Check that the Echo Reply has the same identifier and sequence number
	if len(reply.MessageBody) < 4 {
		t.Fatal("Echo Reply message body too short")
	}
	
	replyIdentifier := uint16(reply.MessageBody[0])<<8 | uint16(reply.MessageBody[1])
	replySequence := uint16(reply.MessageBody[2])<<8 | uint16(reply.MessageBody[3])
	
	if replyIdentifier != identifier {
		t.Errorf("Echo Reply identifier = %v, want %v", replyIdentifier, identifier)
	}
	
	if replySequence != sequenceNumber {
		t.Errorf("Echo Reply sequence number = %v, want %v", replySequence, sequenceNumber)
	}
}

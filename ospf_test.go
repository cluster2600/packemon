package packemon

import (
	"bytes"
	"testing"
)

// TestOSPFBasicFunctionality tests the basic functionality of the OSPF implementation
// OSPF実装の基本的な機能をテストします
func TestOSPFBasicFunctionality(t *testing.T) {
	// Test OSPF packet creation
	// OSPFパケットの作成をテスト
	routerID := uint32(0xC0A80101) // 192.168.1.1
	areaID := uint32(0)            // Backbone area
	messageBody := []byte{0x01, 0x02, 0x03, 0x04}
	
	ospfPacket := NewOSPF(OSPF_TYPE_HELLO, routerID, areaID, messageBody)
	
	// Check OSPF packet fields
	// OSPFパケットのフィールドをチェック
	if ospfPacket.Version != 2 {
		t.Errorf("OSPF version = %d, want %d", ospfPacket.Version, 2)
	}
	
	if ospfPacket.Type != OSPF_TYPE_HELLO {
		t.Errorf("OSPF type = %d, want %d", ospfPacket.Type, OSPF_TYPE_HELLO)
	}
	
	expectedLength := uint16(24 + len(messageBody)) // Header (24) + message body length
	if ospfPacket.PacketLength != expectedLength {
		t.Errorf("OSPF packet length = %d, want %d", ospfPacket.PacketLength, expectedLength)
	}
	
	if ospfPacket.RouterID != routerID {
		t.Errorf("OSPF router ID = %d, want %d", ospfPacket.RouterID, routerID)
	}
	
	if ospfPacket.AreaID != areaID {
		t.Errorf("OSPF area ID = %d, want %d", ospfPacket.AreaID, areaID)
	}
	
	if ospfPacket.AuType != OSPF_AUTH_NONE {
		t.Errorf("OSPF authentication type = %d, want %d", ospfPacket.AuType, OSPF_AUTH_NONE)
	}
	
	if !bytes.Equal(ospfPacket.MessageBody, messageBody) {
		t.Errorf("OSPF message body does not match")
	}
	
	// Check that checksum is non-zero
	// チェックサムがゼロでないことを確認
	if ospfPacket.Checksum == 0 {
		t.Errorf("OSPF checksum is zero, expected non-zero value")
	}
}

// TestOSPFSerialization tests the serialization and deserialization of OSPF packets
// OSPFパケットのシリアル化と逆シリアル化をテストします
func TestOSPFSerialization(t *testing.T) {
	// Create an OSPF Hello packet
	// OSPF Helloパケットを作成
	routerID := uint32(0xC0A80101)      // 192.168.1.1
	areaID := uint32(0)                 // Backbone area
	networkMask := uint32(0xFFFFFF00)   // 255.255.255.0
	helloInterval := uint16(10)         // 10 seconds
	options := uint8(0x02)              // E bit set
	routerPriority := uint8(1)          // Priority 1
	routerDeadInterval := uint32(40)    // 40 seconds
	dr := uint32(0xC0A80101)            // 192.168.1.1
	bdr := uint32(0)                    // 0.0.0.0
	neighbors := []uint32{0xC0A80102}   // 192.168.1.2
	
	ospfHello := NewOSPFHello(routerID, areaID, networkMask, helloInterval, options, routerPriority, routerDeadInterval, dr, bdr, neighbors)
	
	// Serialize the packet
	// パケットをシリアル化
	serialized := ospfHello.Bytes()
	
	// Parse the serialized packet
	// シリアル化されたパケットを解析
	parsed := ParsedOSPF(serialized)
	
	// Check if the parsed packet matches the original
	// 解析されたパケットが元のパケットと一致するかチェック
	if parsed.Version != ospfHello.Version {
		t.Errorf("Parsed OSPF version = %d, want %d", parsed.Version, ospfHello.Version)
	}
	
	if parsed.Type != ospfHello.Type {
		t.Errorf("Parsed OSPF type = %d, want %d", parsed.Type, ospfHello.Type)
	}
	
	if parsed.PacketLength != ospfHello.PacketLength {
		t.Errorf("Parsed OSPF packet length = %d, want %d", parsed.PacketLength, ospfHello.PacketLength)
	}
	
	if parsed.RouterID != ospfHello.RouterID {
		t.Errorf("Parsed OSPF router ID = %d, want %d", parsed.RouterID, ospfHello.RouterID)
	}
	
	if parsed.AreaID != ospfHello.AreaID {
		t.Errorf("Parsed OSPF area ID = %d, want %d", parsed.AreaID, ospfHello.AreaID)
	}
	
	if parsed.Checksum != ospfHello.Checksum {
		t.Errorf("Parsed OSPF checksum = %d, want %d", parsed.Checksum, ospfHello.Checksum)
	}
	
	if parsed.AuType != ospfHello.AuType {
		t.Errorf("Parsed OSPF authentication type = %d, want %d", parsed.AuType, ospfHello.AuType)
	}
	
	// Parse the OSPF Hello packet from the parsed OSPF packet
	// 解析されたOSPFパケットからOSPF Helloパケットを解析
	parsedHello := ParsedOSPFHello(parsed)
	
	// Check OSPF Hello packet fields
	// OSPF Helloパケットのフィールドをチェック
	if parsedHello.NetworkMask != networkMask {
		t.Errorf("Parsed OSPF Hello network mask = %d, want %d", parsedHello.NetworkMask, networkMask)
	}
	
	if parsedHello.HelloInterval != helloInterval {
		t.Errorf("Parsed OSPF Hello interval = %d, want %d", parsedHello.HelloInterval, helloInterval)
	}
	
	if parsedHello.Options != options {
		t.Errorf("Parsed OSPF Hello options = %d, want %d", parsedHello.Options, options)
	}
	
	if parsedHello.RouterPriority != routerPriority {
		t.Errorf("Parsed OSPF Hello router priority = %d, want %d", parsedHello.RouterPriority, routerPriority)
	}
	
	if parsedHello.RouterDeadInterval != routerDeadInterval {
		t.Errorf("Parsed OSPF Hello router dead interval = %d, want %d", parsedHello.RouterDeadInterval, routerDeadInterval)
	}
	
	if parsedHello.DesignatedRouter != dr {
		t.Errorf("Parsed OSPF Hello designated router = %d, want %d", parsedHello.DesignatedRouter, dr)
	}
	
	if parsedHello.BackupDesRouter != bdr {
		t.Errorf("Parsed OSPF Hello backup designated router = %d, want %d", parsedHello.BackupDesRouter, bdr)
	}
	
	if len(parsedHello.Neighbors) != len(neighbors) {
		t.Errorf("Parsed OSPF Hello neighbors length = %d, want %d", len(parsedHello.Neighbors), len(neighbors))
	} else {
		for i, neighbor := range neighbors {
			if parsedHello.Neighbors[i] != neighbor {
				t.Errorf("Parsed OSPF Hello neighbor[%d] = %d, want %d", i, parsedHello.Neighbors[i], neighbor)
			}
		}
	}
}

// TestOSPFChecksum tests the OSPF checksum calculation
// OSPFチェックサム計算をテストします
func TestOSPFChecksum(t *testing.T) {
	// Create an OSPF packet
	// OSPFパケットを作成
	routerID := uint32(0xC0A80101) // 192.168.1.1
	areaID := uint32(0)            // Backbone area
	messageBody := []byte{0x01, 0x02, 0x03, 0x04}
	
	ospfPacket := NewOSPF(OSPF_TYPE_HELLO, routerID, areaID, messageBody)
	
	// Get the original checksum
	// 元のチェックサムを取得
	originalChecksum := ospfPacket.Checksum
	
	// Recalculate the checksum
	// チェックサムを再計算
	recalculatedChecksum := ospfPacket.CalculateChecksum()
	
	// Check if the recalculated checksum matches the original
	// 再計算されたチェックサムが元のチェックサムと一致するかチェック
	if recalculatedChecksum != originalChecksum {
		t.Errorf("Recalculated checksum = %d, want %d", recalculatedChecksum, originalChecksum)
	}
	
	// Modify the packet and check if the checksum changes
	// パケットを変更し、チェックサムが変更されるかチェック
	ospfPacket.RouterID = uint32(0xC0A80102) // 192.168.1.2
	modifiedChecksum := ospfPacket.CalculateChecksum()
	
	if modifiedChecksum == originalChecksum {
		t.Errorf("Modified checksum = %d, should be different from original %d", modifiedChecksum, originalChecksum)
	}
}

// TestOSPFHelloBytes tests the serialization of OSPF Hello packets
// OSPF Helloパケットのシリアル化をテストします
func TestOSPFHelloBytes(t *testing.T) {
	// Create an OSPF Hello packet
	// OSPF Helloパケットを作成
	networkMask := uint32(0xFFFFFF00)   // 255.255.255.0
	helloInterval := uint16(10)         // 10 seconds
	options := uint8(0x02)              // E bit set
	routerPriority := uint8(1)          // Priority 1
	routerDeadInterval := uint32(40)    // 40 seconds
	dr := uint32(0xC0A80101)            // 192.168.1.1
	bdr := uint32(0)                    // 0.0.0.0
	neighbors := []uint32{0xC0A80102}   // 192.168.1.2
	
	hello := &OSPFHello{
		NetworkMask:        networkMask,
		HelloInterval:      helloInterval,
		Options:            options,
		RouterPriority:     routerPriority,
		RouterDeadInterval: routerDeadInterval,
		DesignatedRouter:   dr,
		BackupDesRouter:    bdr,
		Neighbors:          neighbors,
	}
	
	// Serialize the Hello packet
	// Helloパケットをシリアル化
	serialized := hello.Bytes()
	
	// Check the length of the serialized data
	// シリアル化されたデータの長さをチェック
	expectedLength := 20 + (4 * len(neighbors)) // 20 bytes for fixed fields + 4 bytes per neighbor
	if len(serialized) != expectedLength {
		t.Errorf("Serialized OSPF Hello length = %d, want %d", len(serialized), expectedLength)
	}
	
	// Check the network mask field
	// ネットワークマスクフィールドをチェック
	if binary := serialized[0:4]; binary[0] != 0xFF || binary[1] != 0xFF || binary[2] != 0xFF || binary[3] != 0x00 {
		t.Errorf("Serialized network mask = %v, want [255 255 255 0]", binary)
	}
	
	// Check the hello interval field
	// ハロー間隔フィールドをチェック
	if binary := serialized[4:6]; binary[0] != 0x00 || binary[1] != 0x0A {
		t.Errorf("Serialized hello interval = %v, want [0 10]", binary)
	}
	
	// Check the options field
	// オプションフィールドをチェック
	if serialized[6] != options {
		t.Errorf("Serialized options = %d, want %d", serialized[6], options)
	}
	
	// Check the router priority field
	// ルーター優先度フィールドをチェック
	if serialized[7] != routerPriority {
		t.Errorf("Serialized router priority = %d, want %d", serialized[7], routerPriority)
	}
}

// TestOSPFParsingInvalidData tests OSPF parsing with invalid data
// 無効なデータでのOSPF解析をテストします
func TestOSPFParsingInvalidData(t *testing.T) {
	// Test parsing with nil data
	// nilデータでの解析をテスト
	parsed := ParsedOSPF(nil)
	if parsed != nil {
		t.Errorf("ParsedOSPF(nil) = %v, want nil", parsed)
	}
	
	// Test parsing with data that's too short
	// 短すぎるデータでの解析をテスト
	shortData := make([]byte, 23) // OSPF header is 24 bytes
	parsed = ParsedOSPF(shortData)
	if parsed != nil {
		t.Errorf("ParsedOSPF(shortData) = %v, want nil", parsed)
	}
	
	// Test parsing Hello packet with invalid data
	// 無効なデータでのHelloパケット解析をテスト
	invalidHelloData := &OSPF{
		Type: OSPF_TYPE_HELLO,
		MessageBody: []byte{0x01, 0x02}, // Too short for Hello packet
	}
	parsedHello := ParsedOSPFHello(invalidHelloData)
	if parsedHello != nil {
		t.Errorf("ParsedOSPFHello(invalidHelloData) = %v, want nil", parsedHello)
	}
}

// TestOSPFFletcherChecksum tests the Fletcher checksum calculation
// フレッチャーチェックサム計算をテストします
func TestOSPFFletcherChecksum(t *testing.T) {
	// Test data from RFC 1008
	// RFC 1008からのテストデータ
	testData := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x00, 0x00, 0x0E, 0x0F,
	}
	
	// Calculate checksum
	// チェックサムを計算
	checksum := calculateFletcherChecksum(testData)
	
	// Expected checksum from RFC 1008
	// RFC 1008からの期待されるチェックサム
	expectedChecksum := uint16(0xABF5)
	
	if checksum != expectedChecksum {
		t.Errorf("Fletcher checksum = 0x%04X, want 0x%04X", checksum, expectedChecksum)
	}
}

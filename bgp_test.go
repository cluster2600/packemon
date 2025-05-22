package packemon

import (
	"bytes"
	"testing"
)

// TestBGPBasicFunctionality tests the basic functionality of the BGP implementation
// BGP実装の基本的な機能をテストします
func TestBGPBasicFunctionality(t *testing.T) {
	// Test BGP message creation
	// BGPメッセージの作成をテスト
	messageBody := []byte{0x01, 0x02, 0x03, 0x04}
	bgpMessage := NewBGP(BGP_TYPE_UPDATE, messageBody)
	
	// Check BGP message fields
	// BGPメッセージのフィールドをチェック
	if bgpMessage.Type != BGP_TYPE_UPDATE {
		t.Errorf("BGP message type = %d, want %d", bgpMessage.Type, BGP_TYPE_UPDATE)
	}
	
	if !bytes.Equal(bgpMessage.MessageBody, messageBody) {
		t.Errorf("BGP message body does not match")
	}
	
	// Check BGP message length
	// BGPメッセージの長さをチェック
	expectedLength := uint16(19 + len(messageBody)) // Header (19) + message body length
	if bgpMessage.Length != expectedLength {
		t.Errorf("BGP message length = %d, want %d", bgpMessage.Length, expectedLength)
	}
	
	// Check BGP marker
	// BGPマーカーをチェック
	if !bytes.Equal(bgpMessage.Marker, BGP_DEFAULT_MARKER) {
		t.Errorf("BGP marker does not match default marker")
	}
}

// TestBGPSerialization tests the serialization and deserialization of BGP messages
// BGPメッセージのシリアル化と逆シリアル化をテストします
func TestBGPSerialization(t *testing.T) {
	// Create a BGP OPEN message
	// BGP OPENメッセージを作成
	asNumber := uint16(65001)
	holdTime := uint16(180)
	routerID := uint32(0xC0A80101) // 192.168.1.1
	optionalParams := []byte{0x02, 0x06, 0x01, 0x04, 0x00, 0x01, 0x00, 0x01} // Some capabilities
	
	bgpOpen := NewBGPOpen(asNumber, holdTime, routerID, optionalParams)
	
	// Serialize the message
	// メッセージをシリアル化
	serialized := bgpOpen.Bytes()
	
	// Parse the serialized message
	// シリアル化されたメッセージを解析
	parsed := ParsedBGP(serialized)
	
	// Check if the parsed message matches the original
	// 解析されたメッセージが元のメッセージと一致するかチェック
	if parsed.Type != BGP_TYPE_OPEN {
		t.Errorf("Parsed BGP message type = %d, want %d", parsed.Type, BGP_TYPE_OPEN)
	}
	
	if parsed.Length != bgpOpen.Length {
		t.Errorf("Parsed BGP message length = %d, want %d", parsed.Length, bgpOpen.Length)
	}
	
	if !bytes.Equal(parsed.Marker, bgpOpen.Marker) {
		t.Errorf("Parsed BGP marker does not match original")
	}
	
	// Parse the BGP OPEN message from the parsed BGP message
	// 解析されたBGPメッセージからBGP OPENメッセージを解析
	parsedOpen := ParsedBGPOpen(parsed)
	
	// Check BGP OPEN message fields
	// BGP OPENメッセージのフィールドをチェック
	if parsedOpen.Version != 4 {
		t.Errorf("Parsed BGP OPEN version = %d, want %d", parsedOpen.Version, 4)
	}
	
	if parsedOpen.MyAutonomousSystem != asNumber {
		t.Errorf("Parsed BGP OPEN AS number = %d, want %d", parsedOpen.MyAutonomousSystem, asNumber)
	}
	
	if parsedOpen.HoldTime != holdTime {
		t.Errorf("Parsed BGP OPEN hold time = %d, want %d", parsedOpen.HoldTime, holdTime)
	}
	
	if parsedOpen.BGPIdentifier != routerID {
		t.Errorf("Parsed BGP OPEN router ID = %d, want %d", parsedOpen.BGPIdentifier, routerID)
	}
	
	if parsedOpen.OptionalParametersLength != uint8(len(optionalParams)) {
		t.Errorf("Parsed BGP OPEN optional parameters length = %d, want %d", 
			parsedOpen.OptionalParametersLength, len(optionalParams))
	}
	
	if !bytes.Equal(parsedOpen.OptionalParameters, optionalParams) {
		t.Errorf("Parsed BGP OPEN optional parameters do not match original")
	}
}

// TestBGPKeepalive tests the BGP KEEPALIVE message
// BGP KEEPALIVEメッセージをテストします
func TestBGPKeepalive(t *testing.T) {
	// Create a BGP KEEPALIVE message
	// BGP KEEPALIVEメッセージを作成
	keepalive := NewBGPKeepalive()
	
	// Check message type
	// メッセージタイプをチェック
	if keepalive.Type != BGP_TYPE_KEEPALIVE {
		t.Errorf("BGP KEEPALIVE message type = %d, want %d", keepalive.Type, BGP_TYPE_KEEPALIVE)
	}
	
	// Check message body (should be empty)
	// メッセージ本文をチェック（空であるべき）
	if len(keepalive.MessageBody) != 0 {
		t.Errorf("BGP KEEPALIVE message body length = %d, want 0", len(keepalive.MessageBody))
	}
	
	// Check message length (should be 19, header only)
	// メッセージ長をチェック（ヘッダーのみの19であるべき）
	if keepalive.Length != 19 {
		t.Errorf("BGP KEEPALIVE message length = %d, want 19", keepalive.Length)
	}
	
	// Serialize and parse the message
	// メッセージをシリアル化して解析
	serialized := keepalive.Bytes()
	parsed := ParsedBGP(serialized)
	
	// Check if the parsed message matches the original
	// 解析されたメッセージが元のメッセージと一致するかチェック
	if parsed.Type != BGP_TYPE_KEEPALIVE {
		t.Errorf("Parsed BGP KEEPALIVE message type = %d, want %d", parsed.Type, BGP_TYPE_KEEPALIVE)
	}
	
	if len(parsed.MessageBody) != 0 {
		t.Errorf("Parsed BGP KEEPALIVE message body length = %d, want 0", len(parsed.MessageBody))
	}
}

// TestBGPUpdate tests the BGP UPDATE message
// BGP UPDATEメッセージをテストします
func TestBGPUpdate(t *testing.T) {
	// Create test data
	// テストデータを作成
	withdrawnRoutes := []byte{0x18, 0xC0, 0xA8, 0x01} // Withdraw 192.168.1.0/24
	pathAttributes := []byte{
		0x40, 0x01, 0x01, 0x00, // ORIGIN = IGP
		0x40, 0x02, 0x0A, 0x02, 0x01, 0x00, 0x00, 0xFD, 0xE9, 0x00, 0x00, 0xFD, 0xEA, // AS_PATH
	}
	nlri := []byte{0x18, 0xC0, 0xA8, 0x02} // Advertise 192.168.2.0/24
	
	// Create a BGP UPDATE message
	// BGP UPDATEメッセージを作成
	update := NewBGPUpdate(withdrawnRoutes, pathAttributes, nlri)
	
	// Check message type
	// メッセージタイプをチェック
	if update.Type != BGP_TYPE_UPDATE {
		t.Errorf("BGP UPDATE message type = %d, want %d", update.Type, BGP_TYPE_UPDATE)
	}
	
	// Serialize the message
	// メッセージをシリアル化
	serialized := update.Bytes()
	
	// Parse the serialized message
	// シリアル化されたメッセージを解析
	parsed := ParsedBGP(serialized)
	
	// Check if the parsed message matches the original
	// 解析されたメッセージが元のメッセージと一致するかチェック
	if parsed.Type != BGP_TYPE_UPDATE {
		t.Errorf("Parsed BGP UPDATE message type = %d, want %d", parsed.Type, BGP_TYPE_UPDATE)
	}
	
	// Parse the BGP UPDATE message from the parsed BGP message
	// 解析されたBGPメッセージからBGP UPDATEメッセージを解析
	parsedUpdate := ParsedBGPUpdate(parsed)
	
	// Check BGP UPDATE message fields
	// BGP UPDATEメッセージのフィールドをチェック
	if parsedUpdate.WithdrawnRoutesLength != uint16(len(withdrawnRoutes)) {
		t.Errorf("Parsed BGP UPDATE withdrawn routes length = %d, want %d",
			parsedUpdate.WithdrawnRoutesLength, len(withdrawnRoutes))
	}
	
	if !bytes.Equal(parsedUpdate.WithdrawnRoutes, withdrawnRoutes) {
		t.Errorf("Parsed BGP UPDATE withdrawn routes do not match original")
	}
	
	if parsedUpdate.PathAttributesLength != uint16(len(pathAttributes)) {
		t.Errorf("Parsed BGP UPDATE path attributes length = %d, want %d",
			parsedUpdate.PathAttributesLength, len(pathAttributes))
	}
	
	if !bytes.Equal(parsedUpdate.PathAttributes, pathAttributes) {
		t.Errorf("Parsed BGP UPDATE path attributes do not match original")
	}
	
	if !bytes.Equal(parsedUpdate.NetworkLayerReachabilityInfo, nlri) {
		t.Errorf("Parsed BGP UPDATE NLRI does not match original")
	}
}

// TestBGPNotification tests the BGP NOTIFICATION message
// BGP NOTIFICATIONメッセージをテストします
func TestBGPNotification(t *testing.T) {
	// Create test data
	// テストデータを作成
	errorCode := uint8(2)        // Open Message Error
	errorSubcode := uint8(2)     // Bad Peer AS
	data := []byte{0xFD, 0xE9}   // AS 65001
	
	// Create a BGP NOTIFICATION message
	// BGP NOTIFICATIONメッセージを作成
	notification := NewBGPNotification(errorCode, errorSubcode, data)
	
	// Check message type
	// メッセージタイプをチェック
	if notification.Type != BGP_TYPE_NOTIFICATION {
		t.Errorf("BGP NOTIFICATION message type = %d, want %d", notification.Type, BGP_TYPE_NOTIFICATION)
	}
	
	// Serialize the message
	// メッセージをシリアル化
	serialized := notification.Bytes()
	
	// Parse the serialized message
	// シリアル化されたメッセージを解析
	parsed := ParsedBGP(serialized)
	
	// Check if the parsed message matches the original
	// 解析されたメッセージが元のメッセージと一致するかチェック
	if parsed.Type != BGP_TYPE_NOTIFICATION {
		t.Errorf("Parsed BGP NOTIFICATION message type = %d, want %d", parsed.Type, BGP_TYPE_NOTIFICATION)
	}
	
	// Parse the BGP NOTIFICATION message from the parsed BGP message
	// 解析されたBGPメッセージからBGP NOTIFICATIONメッセージを解析
	parsedNotification := ParsedBGPNotification(parsed)
	
	// Check BGP NOTIFICATION message fields
	// BGP NOTIFICATIONメッセージのフィールドをチェック
	if parsedNotification.ErrorCode != errorCode {
		t.Errorf("Parsed BGP NOTIFICATION error code = %d, want %d", parsedNotification.ErrorCode, errorCode)
	}
	
	if parsedNotification.ErrorSubcode != errorSubcode {
		t.Errorf("Parsed BGP NOTIFICATION error subcode = %d, want %d", parsedNotification.ErrorSubcode, errorSubcode)
	}
	
	if !bytes.Equal(parsedNotification.Data, data) {
		t.Errorf("Parsed BGP NOTIFICATION data does not match original")
	}
}

// TestBGPParsingInvalidData tests BGP parsing with invalid data
// 無効なデータでのBGP解析をテストします
func TestBGPParsingInvalidData(t *testing.T) {
	// Test parsing with nil data
	// nilデータでの解析をテスト
	parsed := ParsedBGP(nil)
	if parsed != nil {
		t.Errorf("ParsedBGP(nil) = %v, want nil", parsed)
	}
	
	// Test parsing with data that's too short
	// 短すぎるデータでの解析をテスト
	shortData := make([]byte, 18) // BGP header is 19 bytes
	parsed = ParsedBGP(shortData)
	if parsed != nil {
		t.Errorf("ParsedBGP(shortData) = %v, want nil", parsed)
	}
	
	// Test parsing OPEN message with invalid data
	// 無効なデータでのOPENメッセージ解析をテスト
	invalidOpenData := &BGP{
		Type: BGP_TYPE_OPEN,
		MessageBody: []byte{0x01, 0x02}, // Too short for OPEN message
	}
	parsedOpen := ParsedBGPOpen(invalidOpenData)
	if parsedOpen != nil {
		t.Errorf("ParsedBGPOpen(invalidOpenData) = %v, want nil", parsedOpen)
	}
	
	// Test parsing UPDATE message with invalid data
	// 無効なデータでのUPDATEメッセージ解析をテスト
	invalidUpdateData := &BGP{
		Type: BGP_TYPE_UPDATE,
		MessageBody: []byte{0x01, 0x02}, // Too short for UPDATE message
	}
	parsedUpdate := ParsedBGPUpdate(invalidUpdateData)
	if parsedUpdate != nil {
		t.Errorf("ParsedBGPUpdate(invalidUpdateData) = %v, want nil", parsedUpdate)
	}
	
	// Test parsing NOTIFICATION message with invalid data
	// 無効なデータでのNOTIFICATIONメッセージ解析をテスト
	invalidNotificationData := &BGP{
		Type: BGP_TYPE_NOTIFICATION,
		MessageBody: []byte{0x01}, // Too short for NOTIFICATION message
	}
	parsedNotification := ParsedBGPNotification(invalidNotificationData)
	if parsedNotification != nil {
		t.Errorf("ParsedBGPNotification(invalidNotificationData) = %v, want nil", parsedNotification)
	}
}

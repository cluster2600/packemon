package packemon

import (
	"bytes"
	"encoding/binary"
)

// BGP implements the Border Gateway Protocol as defined in RFC 4271
// BGPはRFC 4271で定義されているボーダーゲートウェイプロトコルを実装します
type BGP struct {
	// BGP Header fields
	// BGPヘッダーフィールド
	Marker                []byte // 16 bytes marker (all 1's for authentication-less BGP) / 16バイトのマーカー（認証なしBGPの場合はすべて1）
	Length                uint16 // Total length of the BGP message / BGPメッセージの総長
	Type                  uint8  // Message type (OPEN, UPDATE, NOTIFICATION, KEEPALIVE) / メッセージタイプ（OPEN、UPDATE、NOTIFICATION、KEEPALIVE）
	MessageBody           []byte // Message-specific data / メッセージ固有のデータ
}

// BGP message types as defined in RFC 4271
// RFC 4271で定義されているBGPメッセージタイプ
const (
	BGP_TYPE_OPEN         = 1
	BGP_TYPE_UPDATE       = 2
	BGP_TYPE_NOTIFICATION = 3
	BGP_TYPE_KEEPALIVE    = 4
)

// BGP OPEN message structure
// BGP OPENメッセージ構造
type BGPOpen struct {
	Version                 uint8  // BGP version number (typically 4) / BGPバージョン番号（通常は4）
	MyAutonomousSystem      uint16 // AS number of the sender / 送信者のAS番号
	HoldTime                uint16 // Proposed hold time in seconds / 提案されたホールドタイム（秒）
	BGPIdentifier           uint32 // BGP identifier (typically router ID) / BGP識別子（通常はルーターID）
	OptionalParametersLength uint8  // Length of optional parameters / オプションパラメータの長さ
	OptionalParameters      []byte // Optional parameters / オプションパラメータ
}

// BGP UPDATE message structure
// BGP UPDATEメッセージ構造
type BGPUpdate struct {
	WithdrawnRoutesLength   uint16 // Length of withdrawn routes / 撤回されたルートの長さ
	WithdrawnRoutes         []byte // Withdrawn routes / 撤回されたルート
	PathAttributesLength    uint16 // Length of path attributes / パス属性の長さ
	PathAttributes          []byte // Path attributes / パス属性
	NetworkLayerReachabilityInfo []byte // NLRI / ネットワーク層到達可能性情報
}

// BGP NOTIFICATION message structure
// BGP NOTIFICATIONメッセージ構造
type BGPNotification struct {
	ErrorCode              uint8  // Error code / エラーコード
	ErrorSubcode           uint8  // Error subcode / エラーサブコード
	Data                   []byte // Error data / エラーデータ
}

// Default BGP marker (all 1's for authentication-less BGP)
// デフォルトのBGPマーカー（認証なしBGPの場合はすべて1）
var BGP_DEFAULT_MARKER = []byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// NewBGP creates a new BGP message with the specified type and message body
// 指定されたタイプとメッセージ本文で新しいBGPメッセージを作成します
func NewBGP(messageType uint8, messageBody []byte) *BGP {
	return &BGP{
		Marker:      BGP_DEFAULT_MARKER,
		Length:      uint16(19 + len(messageBody)), // 19 bytes for header + message body length / ヘッダー19バイト + メッセージ本文の長さ
		Type:        messageType,
		MessageBody: messageBody,
	}
}

// NewBGPOpen creates a new BGP OPEN message
// 新しいBGP OPENメッセージを作成します
func NewBGPOpen(asNumber uint16, holdTime uint16, routerID uint32, optionalParams []byte) *BGP {
	open := &BGPOpen{
		Version:                 4, // BGP version 4 / BGPバージョン4
		MyAutonomousSystem:      asNumber,
		HoldTime:                holdTime,
		BGPIdentifier:           routerID,
		OptionalParametersLength: uint8(len(optionalParams)),
		OptionalParameters:      optionalParams,
	}
	
	// Serialize the OPEN message
	// OPENメッセージをシリアル化
	openBytes := open.Bytes()
	
	// Create a BGP message with the OPEN message as the body
	// OPENメッセージを本文としたBGPメッセージを作成
	return NewBGP(BGP_TYPE_OPEN, openBytes)
}

// NewBGPKeepalive creates a new BGP KEEPALIVE message
// 新しいBGP KEEPALIVEメッセージを作成します
func NewBGPKeepalive() *BGP {
	// KEEPALIVE messages have no message body
	// KEEPALIVEメッセージにはメッセージ本文がありません
	return NewBGP(BGP_TYPE_KEEPALIVE, []byte{})
}

// NewBGPUpdate creates a new BGP UPDATE message
// 新しいBGP UPDATEメッセージを作成します
func NewBGPUpdate(withdrawnRoutes []byte, pathAttributes []byte, nlri []byte) *BGP {
	update := &BGPUpdate{
		WithdrawnRoutesLength:   uint16(len(withdrawnRoutes)),
		WithdrawnRoutes:         withdrawnRoutes,
		PathAttributesLength:    uint16(len(pathAttributes)),
		PathAttributes:          pathAttributes,
		NetworkLayerReachabilityInfo: nlri,
	}
	
	// Serialize the UPDATE message
	// UPDATEメッセージをシリアル化
	updateBytes := update.Bytes()
	
	// Create a BGP message with the UPDATE message as the body
	// UPDATEメッセージを本文としたBGPメッセージを作成
	return NewBGP(BGP_TYPE_UPDATE, updateBytes)
}

// NewBGPNotification creates a new BGP NOTIFICATION message
// 新しいBGP NOTIFICATIONメッセージを作成します
func NewBGPNotification(errorCode uint8, errorSubcode uint8, data []byte) *BGP {
	notification := &BGPNotification{
		ErrorCode:    errorCode,
		ErrorSubcode: errorSubcode,
		Data:         data,
	}
	
	// Serialize the NOTIFICATION message
	// NOTIFICATIONメッセージをシリアル化
	notificationBytes := notification.Bytes()
	
	// Create a BGP message with the NOTIFICATION message as the body
	// NOTIFICATIONメッセージを本文としたBGPメッセージを作成
	return NewBGP(BGP_TYPE_NOTIFICATION, notificationBytes)
}

// Bytes serializes a BGP message into a byte slice
// BGPメッセージをバイトスライスにシリアル化します
func (b *BGP) Bytes() []byte {
	buf := &bytes.Buffer{}
	
	// Write the marker
	// マーカーを書き込む
	buf.Write(b.Marker)
	
	// Write the length
	// 長さを書き込む
	binary.Write(buf, binary.BigEndian, b.Length)
	
	// Write the type
	// タイプを書き込む
	buf.WriteByte(b.Type)
	
	// Write the message body
	// メッセージ本文を書き込む
	buf.Write(b.MessageBody)
	
	return buf.Bytes()
}

// Bytes serializes a BGP OPEN message into a byte slice
// BGP OPENメッセージをバイトスライスにシリアル化します
func (o *BGPOpen) Bytes() []byte {
	buf := &bytes.Buffer{}
	
	// Write the version
	// バージョンを書き込む
	buf.WriteByte(o.Version)
	
	// Write the AS number
	// AS番号を書き込む
	binary.Write(buf, binary.BigEndian, o.MyAutonomousSystem)
	
	// Write the hold time
	// ホールドタイムを書き込む
	binary.Write(buf, binary.BigEndian, o.HoldTime)
	
	// Write the BGP identifier
	// BGP識別子を書き込む
	binary.Write(buf, binary.BigEndian, o.BGPIdentifier)
	
	// Write the optional parameters length
	// オプションパラメータの長さを書き込む
	buf.WriteByte(o.OptionalParametersLength)
	
	// Write the optional parameters
	// オプションパラメータを書き込む
	if o.OptionalParametersLength > 0 {
		buf.Write(o.OptionalParameters)
	}
	
	return buf.Bytes()
}

// Bytes serializes a BGP UPDATE message into a byte slice
// BGP UPDATEメッセージをバイトスライスにシリアル化します
func (u *BGPUpdate) Bytes() []byte {
	buf := &bytes.Buffer{}
	
	// Write the withdrawn routes length
	// 撤回されたルートの長さを書き込む
	binary.Write(buf, binary.BigEndian, u.WithdrawnRoutesLength)
	
	// Write the withdrawn routes
	// 撤回されたルートを書き込む
	if u.WithdrawnRoutesLength > 0 {
		buf.Write(u.WithdrawnRoutes)
	}
	
	// Write the path attributes length
	// パス属性の長さを書き込む
	binary.Write(buf, binary.BigEndian, u.PathAttributesLength)
	
	// Write the path attributes
	// パス属性を書き込む
	if u.PathAttributesLength > 0 {
		buf.Write(u.PathAttributes)
	}
	
	// Write the NLRI
	// NLRIを書き込む
	buf.Write(u.NetworkLayerReachabilityInfo)
	
	return buf.Bytes()
}

// Bytes serializes a BGP NOTIFICATION message into a byte slice
// BGP NOTIFICATIONメッセージをバイトスライスにシリアル化します
func (n *BGPNotification) Bytes() []byte {
	buf := &bytes.Buffer{}
	
	// Write the error code
	// エラーコードを書き込む
	buf.WriteByte(n.ErrorCode)
	
	// Write the error subcode
	// エラーサブコードを書き込む
	buf.WriteByte(n.ErrorSubcode)
	
	// Write the data
	// データを書き込む
	buf.Write(n.Data)
	
	return buf.Bytes()
}

// ParsedBGP parses a BGP message from a byte slice
// バイトスライスからBGPメッセージを解析します
func ParsedBGP(data []byte) *BGP {
	if len(data) < 19 { // Minimum BGP message size is 19 bytes / BGPメッセージの最小サイズは19バイト
		return nil
	}
	
	return &BGP{
		Marker:      data[0:16],
		Length:      binary.BigEndian.Uint16(data[16:18]),
		Type:        data[18],
		MessageBody: data[19:],
	}
}

// ParsedBGPOpen parses a BGP OPEN message from a BGP message
// BGPメッセージからBGP OPENメッセージを解析します
func ParsedBGPOpen(bgp *BGP) *BGPOpen {
	if bgp == nil || bgp.Type != BGP_TYPE_OPEN || len(bgp.MessageBody) < 10 {
		return nil
	}
	
	optParamLen := bgp.MessageBody[9]
	var optParams []byte
	
	if optParamLen > 0 && len(bgp.MessageBody) >= 10+int(optParamLen) {
		optParams = bgp.MessageBody[10:10+optParamLen]
	}
	
	return &BGPOpen{
		Version:                 bgp.MessageBody[0],
		MyAutonomousSystem:      binary.BigEndian.Uint16(bgp.MessageBody[1:3]),
		HoldTime:                binary.BigEndian.Uint16(bgp.MessageBody[3:5]),
		BGPIdentifier:           binary.BigEndian.Uint32(bgp.MessageBody[5:9]),
		OptionalParametersLength: optParamLen,
		OptionalParameters:      optParams,
	}
}

// ParsedBGPUpdate parses a BGP UPDATE message from a BGP message
// BGPメッセージからBGP UPDATEメッセージを解析します
func ParsedBGPUpdate(bgp *BGP) *BGPUpdate {
	if bgp == nil || bgp.Type != BGP_TYPE_UPDATE || len(bgp.MessageBody) < 4 {
		return nil
	}
	
	withdrawnRoutesLen := binary.BigEndian.Uint16(bgp.MessageBody[0:2])
	
	if len(bgp.MessageBody) < 2+int(withdrawnRoutesLen)+2 {
		return nil
	}
	
	withdrawnRoutes := bgp.MessageBody[2:2+withdrawnRoutesLen]
	
	pathAttrLenPos := 2 + withdrawnRoutesLen
	pathAttrLen := binary.BigEndian.Uint16(bgp.MessageBody[pathAttrLenPos:pathAttrLenPos+2])
	
	if len(bgp.MessageBody) < 2+int(withdrawnRoutesLen)+2+int(pathAttrLen) {
		return nil
	}
	
	pathAttr := bgp.MessageBody[pathAttrLenPos+2:pathAttrLenPos+2+pathAttrLen]
	nlri := bgp.MessageBody[pathAttrLenPos+2+pathAttrLen:]
	
	return &BGPUpdate{
		WithdrawnRoutesLength:   withdrawnRoutesLen,
		WithdrawnRoutes:         withdrawnRoutes,
		PathAttributesLength:    pathAttrLen,
		PathAttributes:          pathAttr,
		NetworkLayerReachabilityInfo: nlri,
	}
}

// ParsedBGPNotification parses a BGP NOTIFICATION message from a BGP message
// BGPメッセージからBGP NOTIFICATIONメッセージを解析します
func ParsedBGPNotification(bgp *BGP) *BGPNotification {
	if bgp == nil || bgp.Type != BGP_TYPE_NOTIFICATION || len(bgp.MessageBody) < 2 {
		return nil
	}
	
	return &BGPNotification{
		ErrorCode:    bgp.MessageBody[0],
		ErrorSubcode: bgp.MessageBody[1],
		Data:         bgp.MessageBody[2:],
	}
}

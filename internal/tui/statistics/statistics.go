package statistics

import (
	"net"
	"sort"
	"sync"
	"time"

	"github.com/ddddddO/packemon"
)

// Statistics represents packet statistics
// Statisticsはパケット統計を表します
type Statistics struct {
	// General statistics
	// 一般統計
	startTime      time.Time
	totalPackets   int
	totalBytes     int64
	
	// Protocol statistics
	// プロトコル統計
	protocolCounts map[string]int
	
	// IP statistics
	// IP統計
	sourceIPs      map[string]int
	destIPs        map[string]int
	
	// Packet rate statistics
	// パケットレート統計
	packetCounts   []int
	lastCountTime  time.Time
	currentCount   int
	
	// Mutex for thread safety
	// スレッドセーフのためのミューテックス
	mu             sync.Mutex
}

// IPCount represents an IP address and its packet count
// IPCountはIPアドレスとそのパケット数を表します
type IPCount struct {
	IP    string
	Count int
}

// NewStatistics creates a new statistics object
// 新しい統計オブジェクトを作成します
func NewStatistics() *Statistics {
	return &Statistics{
		startTime:      time.Now(),
		protocolCounts: make(map[string]int),
		sourceIPs:      make(map[string]int),
		destIPs:        make(map[string]int),
		packetCounts:   make([]int, 60), // Store 60 seconds of history / 60秒間の履歴を保存
		lastCountTime:  time.Now(),
	}
}

// ProcessPacket processes a packet for statistics
// 統計のためにパケットを処理します
func (s *Statistics) ProcessPacket(passive *packemon.Passive) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Update total packet count and size
	// 総パケット数とサイズを更新
	s.totalPackets++
	
	// Calculate packet size
	// パケットサイズを計算
	packetSize := s.calculatePacketSize(passive)
	s.totalBytes += int64(packetSize)
	
	// Update protocol statistics
	// プロトコル統計を更新
	s.updateProtocolStats(passive)
	
	// Update IP statistics
	// IP統計を更新
	s.updateIPStats(passive)
	
	// Update packet rate statistics
	// パケットレート統計を更新
	s.updatePacketRateStats()
}

// calculatePacketSize calculates the size of a packet
// パケットのサイズを計算します
func (s *Statistics) calculatePacketSize(passive *packemon.Passive) int {
	size := 0
	
	// Add Ethernet frame size if available
	// イーサネットフレームサイズが利用可能な場合は追加
	if passive.EthernetFrame != nil {
		size = len(passive.EthernetFrame.Bytes())
	} else {
		// Otherwise estimate size based on available layers
		// それ以外の場合は、利用可能なレイヤーに基づいてサイズを推定
		
		// Add IPv4 size
		// IPv4サイズを追加
		if passive.IPv4 != nil {
			size += len(passive.IPv4.Bytes())
		}
		
		// Add IPv6 size
		// IPv6サイズを追加
		if passive.IPv6 != nil {
			size += len(passive.IPv6.Bytes())
		}
		
		// Add TCP size
		// TCPサイズを追加
		if passive.TCP != nil {
			size += len(passive.TCP.Bytes())
		}
		
		// Add UDP size
		// UDPサイズを追加
		if passive.UDP != nil {
			size += len(passive.UDP.Bytes())
		}
		
		// Add ICMP size
		// ICMPサイズを追加
		if passive.ICMP != nil {
			size += len(passive.ICMP.Bytes())
		}
		
		// Add ICMPv6 size
		// ICMPv6サイズを追加
		if passive.ICMPv6 != nil {
			size += len(passive.ICMPv6.Bytes())
		}
	}
	
	return size
}

// updateProtocolStats updates protocol statistics
// プロトコル統計を更新します
func (s *Statistics) updateProtocolStats(passive *packemon.Passive) {
	// Update Ethernet count
	// イーサネット数を更新
	if passive.EthernetFrame != nil {
		s.protocolCounts["Ethernet"]++
	}
	
	// Update IPv4 count
	// IPv4数を更新
	if passive.IPv4 != nil {
		s.protocolCounts["IPv4"]++
	}
	
	// Update IPv6 count
	// IPv6数を更新
	if passive.IPv6 != nil {
		s.protocolCounts["IPv6"]++
	}
	
	// Update TCP count
	// TCP数を更新
	if passive.TCP != nil {
		s.protocolCounts["TCP"]++
	}
	
	// Update UDP count
	// UDP数を更新
	if passive.UDP != nil {
		s.protocolCounts["UDP"]++
	}
	
	// Update ICMP count
	// ICMP数を更新
	if passive.ICMP != nil {
		s.protocolCounts["ICMP"]++
	}
	
	// Update ICMPv6 count
	// ICMPv6数を更新
	if passive.ICMPv6 != nil {
		s.protocolCounts["ICMPv6"]++
	}
	
	// Update DNS count
	// DNS数を更新
	if passive.DNS != nil {
		s.protocolCounts["DNS"]++
	}
	
	// Update HTTP count
	// HTTP数を更新
	if passive.HTTP != nil {
		s.protocolCounts["HTTP"]++
	}
	
	// Update TLS count
	// TLS数を更新
	if passive.TLSClientHello != nil || passive.TLSServerHello != nil {
		s.protocolCounts["TLS"]++
	}
	
	// Update ARP count
	// ARP数を更新
	if passive.ARP != nil {
		s.protocolCounts["ARP"]++
	}
	
	// Update BGP count
	// BGP数を更新
	if passive.BGP != nil {
		s.protocolCounts["BGP"]++
	}
	
	// Update OSPF count
	// OSPF数を更新
	if passive.OSPF != nil {
		s.protocolCounts["OSPF"]++
	}
}

// updateIPStats updates IP statistics
// IP統計を更新します
func (s *Statistics) updateIPStats(passive *packemon.Passive) {
	// Get source and destination IP addresses
	// 送信元と宛先のIPアドレスを取得
	var srcIP, dstIP net.IP
	
	if passive.IPv4 != nil {
		srcIP = passive.IPv4.SrcAddr
		dstIP = passive.IPv4.DstAddr
	} else if passive.IPv6 != nil {
		srcIP = passive.IPv6.SrcAddr
		dstIP = passive.IPv6.DstAddr
	}
	
	// Update source IP count
	// 送信元IP数を更新
	if srcIP != nil {
		s.sourceIPs[srcIP.String()]++
	}
	
	// Update destination IP count
	// 宛先IP数を更新
	if dstIP != nil {
		s.destIPs[dstIP.String()]++
	}
}

// updatePacketRateStats updates packet rate statistics
// パケットレート統計を更新します
func (s *Statistics) updatePacketRateStats() {
	// Increment current count
	// 現在のカウントをインクリメント
	s.currentCount++
	
	// Check if a second has passed
	// 1秒が経過したかどうかを確認
	now := time.Now()
	if now.Sub(s.lastCountTime) >= time.Second {
		// Shift counts to the left
		// カウントを左にシフト
		copy(s.packetCounts[0:], s.packetCounts[1:])
		
		// Add current count to the end
		// 現在のカウントを最後に追加
		s.packetCounts[len(s.packetCounts)-1] = s.currentCount
		
		// Reset current count and update last count time
		// 現在のカウントをリセットし、最後のカウント時間を更新
		s.currentCount = 0
		s.lastCountTime = now
	}
}

// TotalPackets returns the total number of packets
// パケットの総数を返します
func (s *Statistics) TotalPackets() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.totalPackets
}

// TotalBytes returns the total number of bytes
// バイトの総数を返します
func (s *Statistics) TotalBytes() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.totalBytes
}

// AveragePacketSize returns the average packet size
// 平均パケットサイズを返します
func (s *Statistics) AveragePacketSize() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.totalPackets == 0 {
		return 0
	}
	
	return float64(s.totalBytes) / float64(s.totalPackets)
}

// PacketRate returns the current packet rate
// 現在のパケットレートを返します
func (s *Statistics) PacketRate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Calculate packets per second based on total time
	// 総時間に基づいて1秒あたりのパケット数を計算
	duration := time.Since(s.startTime).Seconds()
	if duration <= 0 {
		return 0
	}
	
	return float64(s.totalPackets) / duration
}

// MonitoringTime returns the total monitoring time
// 総モニタリング時間を返します
func (s *Statistics) MonitoringTime() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return time.Since(s.startTime)
}

// ProtocolDistribution returns the protocol distribution
// プロトコル分布を返します
func (s *Statistics) ProtocolDistribution() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create a copy of the protocol counts
	// プロトコル数のコピーを作成
	counts := make(map[string]int)
	for proto, count := range s.protocolCounts {
		counts[proto] = count
	}
	
	return counts
}

// PacketRateHistory returns the packet rate history
// パケットレート履歴を返します
func (s *Statistics) PacketRateHistory() []float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Convert packet counts to rates
	// パケット数をレートに変換
	rates := make([]float64, len(s.packetCounts))
	for i, count := range s.packetCounts {
		rates[i] = float64(count)
	}
	
	return rates
}

// TopSourceIPs returns the top source IPs
// トップ送信元IPを返します
func (s *Statistics) TopSourceIPs(n int) []IPCount {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.topIPs(s.sourceIPs, n)
}

// TopDestinationIPs returns the top destination IPs
// トップ宛先IPを返します
func (s *Statistics) TopDestinationIPs(n int) []IPCount {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.topIPs(s.destIPs, n)
}

// topIPs returns the top n IPs from the given map
// 指定されたマップからトップnのIPを返します
func (s *Statistics) topIPs(ips map[string]int, n int) []IPCount {
	// Create a slice of IPCount
	// IPCountのスライスを作成
	ipCounts := make([]IPCount, 0, len(ips))
	for ip, count := range ips {
		ipCounts = append(ipCounts, IPCount{IP: ip, Count: count})
	}
	
	// Sort by count in descending order
	// カウントの降順でソート
	sort.Slice(ipCounts, func(i, j int) bool {
		return ipCounts[i].Count > ipCounts[j].Count
	})
	
	// Return top n
	// トップnを返す
	if len(ipCounts) > n {
		return ipCounts[:n]
	}
	
	return ipCounts
}

// Reset resets all statistics
// すべての統計をリセットします
func (s *Statistics) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.startTime = time.Now()
	s.totalPackets = 0
	s.totalBytes = 0
	s.protocolCounts = make(map[string]int)
	s.sourceIPs = make(map[string]int)
	s.destIPs = make(map[string]int)
	s.packetCounts = make([]int, 60)
	s.lastCountTime = time.Now()
	s.currentCount = 0
}

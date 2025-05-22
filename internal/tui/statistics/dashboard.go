package statistics

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ddddddO/packemon"
)

// Dashboard represents a statistics dashboard for packet monitoring
// Dashboardはパケットモニタリングの統計ダッシュボードを表します
type Dashboard struct {
	app            *tview.Application
	flex           *tview.Flex
	packetCountBox *tview.TextView
	protocolChart  *tview.TextView
	timelineChart  *tview.TextView
	topTalkers     *tview.TextView
	
	// Statistics data
	// 統計データ
	stats          *Statistics
	
	// Mutex for thread safety
	// スレッドセーフのためのミューテックス
	mu             sync.Mutex
	
	// Update ticker
	// 更新用ティッカー
	ticker         *time.Ticker
	done           chan bool
}

// NewDashboard creates a new statistics dashboard
// 新しい統計ダッシュボードを作成します
func NewDashboard(app *tview.Application) *Dashboard {
	d := &Dashboard{
		app:   app,
		stats: NewStatistics(),
		done:  make(chan bool),
	}
	
	// Initialize UI components
	// UIコンポーネントを初期化
	d.initUI()
	
	// Start update ticker (update every second)
	// 更新用ティッカーを開始（1秒ごとに更新）
	d.ticker = time.NewTicker(1 * time.Second)
	go d.updateLoop()
	
	return d
}

// initUI initializes the UI components
// UIコンポーネントを初期化します
func (d *Dashboard) initUI() {
	// Packet count box
	// パケット数ボックス
	d.packetCountBox = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTitle("Packet Statistics").
		SetBorder(true)
	
	// Protocol distribution chart
	// プロトコル分布チャート
	d.protocolChart = tview.NewTextView().
		SetDynamicColors(true).
		SetTitle("Protocol Distribution").
		SetBorder(true)
	
	// Timeline chart
	// タイムラインチャート
	d.timelineChart = tview.NewTextView().
		SetDynamicColors(true).
		SetTitle("Packet Rate (packets/sec)").
		SetBorder(true)
	
	// Top talkers
	// トップトーカー
	d.topTalkers = tview.NewTextView().
		SetDynamicColors(true).
		SetTitle("Top Talkers").
		SetBorder(true)
	
	// Create layout
	// レイアウトを作成
	topRow := tview.NewFlex().
		AddItem(d.packetCountBox, 0, 1, false).
		AddItem(d.protocolChart, 0, 2, false)
	
	bottomRow := tview.NewFlex().
		AddItem(d.timelineChart, 0, 2, false).
		AddItem(d.topTalkers, 0, 1, false)
	
	d.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topRow, 0, 1, false).
		AddItem(bottomRow, 0, 1, false)
}

// updateLoop updates the dashboard periodically
// ダッシュボードを定期的に更新します
func (d *Dashboard) updateLoop() {
	for {
		select {
		case <-d.ticker.C:
			d.updateUI()
		case <-d.done:
			return
		}
	}
}

// updateUI updates the UI components with current statistics
// 現在の統計情報でUIコンポーネントを更新します
func (d *Dashboard) updateUI() {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.app.QueueUpdateDraw(func() {
		// Update packet count box
		// パケット数ボックスを更新
		d.updatePacketCountBox()
		
		// Update protocol distribution chart
		// プロトコル分布チャートを更新
		d.updateProtocolChart()
		
		// Update timeline chart
		// タイムラインチャートを更新
		d.updateTimelineChart()
		
		// Update top talkers
		// トップトーカーを更新
		d.updateTopTalkers()
	})
}

// updatePacketCountBox updates the packet count box
// パケット数ボックスを更新します
func (d *Dashboard) updatePacketCountBox() {
	d.packetCountBox.Clear()
	
	totalPackets := d.stats.TotalPackets()
	avgSize := d.stats.AveragePacketSize()
	packetRate := d.stats.PacketRate()
	
	fmt.Fprintf(d.packetCountBox, "[yellow]Total Packets:[white] %d\n", totalPackets)
	fmt.Fprintf(d.packetCountBox, "[yellow]Average Size:[white] %.2f bytes\n", avgSize)
	fmt.Fprintf(d.packetCountBox, "[yellow]Packet Rate:[white] %.2f pps\n", packetRate)
	fmt.Fprintf(d.packetCountBox, "[yellow]Monitoring Time:[white] %s\n", d.stats.MonitoringTime().String())
}

// updateProtocolChart updates the protocol distribution chart
// プロトコル分布チャートを更新します
func (d *Dashboard) updateProtocolChart() {
	d.protocolChart.Clear()
	
	// Get protocol distribution
	// プロトコル分布を取得
	protocols := d.stats.ProtocolDistribution()
	
	// Find the maximum count for scaling
	// スケーリングのための最大カウントを見つける
	maxCount := 0
	for _, count := range protocols {
		if count > maxCount {
			maxCount = count
		}
	}
	
	// Draw the chart
	// チャートを描画
	for proto, count := range protocols {
		// Calculate bar length (max 40 characters)
		// バーの長さを計算（最大40文字）
		barLength := 0
		if maxCount > 0 {
			barLength = count * 40 / maxCount
		}
		
		// Create the bar
		// バーを作成
		bar := ""
		for i := 0; i < barLength; i++ {
			bar += "█"
		}
		
		// Calculate percentage
		// パーセンテージを計算
		percentage := 0.0
		if d.stats.TotalPackets() > 0 {
			percentage = float64(count) * 100.0 / float64(d.stats.TotalPackets())
		}
		
		// Print the bar
		// バーを表示
		fmt.Fprintf(d.protocolChart, "[yellow]%-8s[green]%s [white]%d [blue](%.1f%%)\n", proto, bar, count, percentage)
	}
}

// updateTimelineChart updates the timeline chart
// タイムラインチャートを更新します
func (d *Dashboard) updateTimelineChart() {
	d.timelineChart.Clear()
	
	// Get packet rate history
	// パケットレート履歴を取得
	history := d.stats.PacketRateHistory()
	
	// Find the maximum rate for scaling
	// スケーリングのための最大レートを見つける
	maxRate := 0.0
	for _, rate := range history {
		if rate > maxRate {
			maxRate = rate
		}
	}
	
	// Ensure we have a non-zero max for scaling
	// スケーリングのために非ゼロの最大値を確保
	if maxRate < 1.0 {
		maxRate = 1.0
	}
	
	// Draw the chart
	// チャートを描画
	for i, rate := range history {
		// Calculate bar length (max 60 characters)
		// バーの長さを計算（最大60文字）
		barLength := int(rate * 60.0 / maxRate)
		
		// Create the bar
		// バーを作成
		bar := ""
		for j := 0; j < barLength; j++ {
			bar += "█"
		}
		
		// Print the bar with timestamp
		// タイムスタンプ付きでバーを表示
		timeAgo := len(history) - i - 1
		fmt.Fprintf(d.timelineChart, "[yellow]%2ds ago:[blue]%s [white]%.2f pps\n", timeAgo, bar, rate)
	}
}

// updateTopTalkers updates the top talkers display
// トップトーカー表示を更新します
func (d *Dashboard) updateTopTalkers() {
	d.topTalkers.Clear()
	
	// Get top source IPs
	// トップ送信元IPを取得
	srcIPs := d.stats.TopSourceIPs(5)
	
	// Print top source IPs
	// トップ送信元IPを表示
	fmt.Fprintf(d.topTalkers, "[yellow]Top Source IPs:\n")
	for i, entry := range srcIPs {
		fmt.Fprintf(d.topTalkers, "[white]%d. [green]%s [white]- %d packets\n", i+1, entry.IP, entry.Count)
	}
	
	fmt.Fprintf(d.topTalkers, "\n")
	
	// Get top destination IPs
	// トップ宛先IPを取得
	dstIPs := d.stats.TopDestinationIPs(5)
	
	// Print top destination IPs
	// トップ宛先IPを表示
	fmt.Fprintf(d.topTalkers, "[yellow]Top Destination IPs:\n")
	for i, entry := range dstIPs {
		fmt.Fprintf(d.topTalkers, "[white]%d. [green]%s [white]- %d packets\n", i+1, entry.IP, entry.Count)
	}
}

// ProcessPacket processes a packet for statistics
// 統計のためにパケットを処理します
func (d *Dashboard) ProcessPacket(passive *packemon.Passive) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.stats.ProcessPacket(passive)
}

// GetView returns the main view of the dashboard
// ダッシュボードのメインビューを返します
func (d *Dashboard) GetView() tview.Primitive {
	return d.flex
}

// Stop stops the dashboard updates
// ダッシュボードの更新を停止します
func (d *Dashboard) Stop() {
	d.ticker.Stop()
	d.done <- true
}

// HandleKey handles key events
// キーイベントを処理します
func (d *Dashboard) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	// Handle key events here
	// ここでキーイベントを処理
	
	// For now, just pass the event through
	// 今のところ、イベントをそのまま通過させる
	return event
}

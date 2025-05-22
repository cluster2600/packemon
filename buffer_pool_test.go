package packemon

import (
	"bytes"
	"testing"
)

// TestBufferPool tests the buffer pool functionality
// バッファプールの機能をテストします
func TestBufferPool(t *testing.T) {
	// Create a new buffer pool
	// 新しいバッファプールを作成
	pool := NewBufferPool()
	
	// Get a buffer from the pool
	// プールからバッファを取得
	buf := pool.Get()
	
	// Check that the buffer is empty
	// バッファが空であることを確認
	if buf.Len() != 0 {
		t.Errorf("Buffer length = %d, want 0", buf.Len())
	}
	
	// Write some data to the buffer
	// バッファにデータを書き込む
	testData := "Hello, world!"
	buf.WriteString(testData)
	
	// Check that the buffer contains the data
	// バッファにデータが含まれていることを確認
	if buf.String() != testData {
		t.Errorf("Buffer content = %q, want %q", buf.String(), testData)
	}
	
	// Put the buffer back in the pool
	// バッファをプールに戻す
	pool.Put(buf)
	
	// Get another buffer from the pool
	// プールから別のバッファを取得
	buf2 := pool.Get()
	
	// Check that the buffer is empty (should be reset)
	// バッファが空であることを確認（リセットされているはず）
	if buf2.Len() != 0 {
		t.Errorf("Buffer length after reset = %d, want 0", buf2.Len())
	}
	
	// Check that we got the same buffer back (this is implementation-dependent)
	// 同じバッファが返されたことを確認（これは実装依存）
	// We can't directly check if buf == buf2 because sync.Pool might return a different object
	// sync.Poolは異なるオブジェクトを返す可能性があるため、buf == buf2を直接チェックすることはできない
}

// TestGlobalBufferPool tests the global buffer pool functions
// グローバルバッファプール関数をテストします
func TestGlobalBufferPool(t *testing.T) {
	// Get a buffer from the global pool
	// グローバルプールからバッファを取得
	buf := GetBuffer()
	
	// Check that the buffer is empty
	// バッファが空であることを確認
	if buf.Len() != 0 {
		t.Errorf("Buffer length = %d, want 0", buf.Len())
	}
	
	// Write some data to the buffer
	// バッファにデータを書き込む
	testData := "Hello, global pool!"
	buf.WriteString(testData)
	
	// Check that the buffer contains the data
	// バッファにデータが含まれていることを確認
	if buf.String() != testData {
		t.Errorf("Buffer content = %q, want %q", buf.String(), testData)
	}
	
	// Put the buffer back in the pool
	// バッファをプールに戻す
	PutBuffer(buf)
	
	// Get another buffer from the pool
	// プールから別のバッファを取得
	buf2 := GetBuffer()
	
	// Check that the buffer is empty (should be reset)
	// バッファが空であることを確認（リセットされているはず）
	if buf2.Len() != 0 {
		t.Errorf("Buffer length after reset = %d, want 0", buf2.Len())
	}
}

// TestBytesPool tests the byte slice pool functionality
// バイトスライスプールの機能をテストします
func TestBytesPool(t *testing.T) {
	// Create a new bytes pool with size 10
	// サイズ10の新しいバイトプールを作成
	pool := NewBytesPool(10)
	
	// Get a byte slice from the pool
	// プールからバイトスライスを取得
	buf := pool.Get()
	
	// Check that the byte slice has the correct size
	// バイトスライスが正しいサイズであることを確認
	if len(buf) != 10 {
		t.Errorf("Byte slice length = %d, want 10", len(buf))
	}
	
	// Check that the byte slice is zeroed
	// バイトスライスがゼロ化されていることを確認
	for i, b := range buf {
		if b != 0 {
			t.Errorf("Byte slice[%d] = %d, want 0", i, b)
		}
	}
	
	// Fill the byte slice with data
	// バイトスライスにデータを入力
	for i := range buf {
		buf[i] = byte(i + 1)
	}
	
	// Put the byte slice back in the pool
	// バイトスライスをプールに戻す
	pool.Put(buf)
	
	// Get another byte slice from the pool
	// プールから別のバイトスライスを取得
	buf2 := pool.Get()
	
	// Check that the byte slice is zeroed (should be cleared for security)
	// バイトスライスがゼロ化されていることを確認（セキュリティのためにクリアされているはず）
	for i, b := range buf2 {
		if b != 0 {
			t.Errorf("Byte slice[%d] after reset = %d, want 0", i, b)
		}
	}
}

// TestGlobalBytesPools tests the global byte slice pool functions
// グローバルバイトスライスプール関数をテストします
func TestGlobalBytesPools(t *testing.T) {
	// Test small bytes pool
	// 小さいバイトプールをテスト
	smallBuf := GetSmallBytes()
	if len(smallBuf) != SmallPacketSize {
		t.Errorf("Small byte slice length = %d, want %d", len(smallBuf), SmallPacketSize)
	}
	PutSmallBytes(smallBuf)
	
	// Test medium bytes pool
	// 中サイズのバイトプールをテスト
	mediumBuf := GetMediumBytes()
	if len(mediumBuf) != MediumPacketSize {
		t.Errorf("Medium byte slice length = %d, want %d", len(mediumBuf), MediumPacketSize)
	}
	PutMediumBytes(mediumBuf)
	
	// Test large bytes pool
	// 大きいバイトプールをテスト
	largeBuf := GetLargeBytes()
	if len(largeBuf) != LargePacketSize {
		t.Errorf("Large byte slice length = %d, want %d", len(largeBuf), LargePacketSize)
	}
	PutLargeBytes(largeBuf)
}

// TestGetBytes tests the GetBytes function
// GetBytes関数をテストします
func TestGetBytes(t *testing.T) {
	// Test getting small bytes
	// 小さいバイトの取得をテスト
	smallBuf := GetBytes(100)
	if len(smallBuf) < 100 {
		t.Errorf("Small byte slice length = %d, want >= 100", len(smallBuf))
	}
	
	// Test getting medium bytes
	// 中サイズのバイトの取得をテスト
	mediumBuf := GetBytes(1000)
	if len(mediumBuf) < 1000 {
		t.Errorf("Medium byte slice length = %d, want >= 1000", len(mediumBuf))
	}
	
	// Test getting large bytes
	// 大きいバイトの取得をテスト
	largeBuf := GetBytes(5000)
	if len(largeBuf) < 5000 {
		t.Errorf("Large byte slice length = %d, want >= 5000", len(largeBuf))
	}
}

// TestPutBytes tests the PutBytes function
// PutBytes関数をテストします
func TestPutBytes(t *testing.T) {
	// Create byte slices of different sizes
	// 異なるサイズのバイトスライスを作成
	smallBuf := make([]byte, 100)
	mediumBuf := make([]byte, 1000)
	largeBuf := make([]byte, 5000)
	tooLargeBuf := make([]byte, 10000)
	
	// Put the byte slices back in the pool
	// バイトスライスをプールに戻す
	PutBytes(smallBuf)
	PutBytes(mediumBuf)
	PutBytes(largeBuf)
	PutBytes(tooLargeBuf) // This should not cause any errors
	
	// We can't directly test if the byte slices were put back in the pool
	// バイトスライスがプールに戻されたかどうかを直接テストすることはできない
	// But we can check that the function doesn't panic
	// しかし、関数がパニックを起こさないことを確認できる
}

// TestBufferPoolPerformance tests the performance of the buffer pool
// バッファプールのパフォーマンスをテストします
func TestBufferPoolPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	// Create a new buffer pool
	// 新しいバッファプールを作成
	pool := NewBufferPool()
	
	// Number of iterations
	// 繰り返し回数
	iterations := 100000
	
	// Test getting and putting buffers
	// バッファの取得と返却をテスト
	for i := 0; i < iterations; i++ {
		buf := pool.Get()
		buf.WriteString("Test data")
		pool.Put(buf)
	}
	
	// Test creating new buffers each time (for comparison)
	// 比較のために毎回新しいバッファを作成するテスト
	for i := 0; i < iterations; i++ {
		buf := new(bytes.Buffer)
		buf.WriteString("Test data")
		_ = buf
	}
	
	// No assertions here, this is just to measure performance
	// ここにはアサーションはなく、パフォーマンスを測定するだけ
}

// TestBytesPoolPerformance tests the performance of the bytes pool
// バイトプールのパフォーマンスをテストします
func TestBytesPoolPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	// Create a new bytes pool
	// 新しいバイトプールを作成
	pool := NewBytesPool(1500)
	
	// Number of iterations
	// 繰り返し回数
	iterations := 100000
	
	// Test getting and putting byte slices
	// バイトスライスの取得と返却をテスト
	for i := 0; i < iterations; i++ {
		buf := pool.Get()
		buf[0] = 1
		pool.Put(buf)
	}
	
	// Test creating new byte slices each time (for comparison)
	// 比較のために毎回新しいバイトスライスを作成するテスト
	for i := 0; i < iterations; i++ {
		buf := make([]byte, 1500)
		buf[0] = 1
		_ = buf
	}
	
	// No assertions here, this is just to measure performance
	// ここにはアサーションはなく、パフォーマンスを測定するだけ
}

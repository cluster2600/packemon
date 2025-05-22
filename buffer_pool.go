package packemon

import (
	"bytes"
	"sync"
)

// BufferPool implements a pool of bytes.Buffer objects to reduce GC pressure
// BufferPoolはGCの負荷を軽減するためのbytes.Bufferオブジェクトのプールを実装します
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new buffer pool
// 新しいバッファプールを作成します
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get retrieves a buffer from the pool
// プールからバッファを取得します
func (p *BufferPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset() // Ensure the buffer is empty / バッファが空であることを確認
	return buf
}

// Put returns a buffer to the pool
// バッファをプールに返します
func (p *BufferPool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

// Global buffer pool instance
// グローバルバッファプールインスタンス
var globalBufferPool = NewBufferPool()

// GetBuffer retrieves a buffer from the global pool
// グローバルプールからバッファを取得します
func GetBuffer() *bytes.Buffer {
	return globalBufferPool.Get()
}

// PutBuffer returns a buffer to the global pool
// バッファをグローバルプールに返します
func PutBuffer(buf *bytes.Buffer) {
	globalBufferPool.Put(buf)
}

// BytesPool implements a pool of byte slices to reduce GC pressure
// BytesPoolはGCの負荷を軽減するためのバイトスライスのプールを実装します
type BytesPool struct {
	pool sync.Pool
	size int
}

// NewBytesPool creates a new bytes pool with the specified size
// 指定されたサイズの新しいバイトプールを作成します
func NewBytesPool(size int) *BytesPool {
	return &BytesPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
		size: size,
	}
}

// Get retrieves a byte slice from the pool
// プールからバイトスライスを取得します
func (p *BytesPool) Get() []byte {
	buf := p.pool.Get().([]byte)
	// Clear the buffer for security reasons
	// セキュリティ上の理由でバッファをクリア
	for i := range buf {
		buf[i] = 0
	}
	return buf
}

// Put returns a byte slice to the pool
// バイトスライスをプールに返します
func (p *BytesPool) Put(buf []byte) {
	if cap(buf) >= p.size {
		p.pool.Put(buf[:p.size])
	}
	// If the buffer is smaller than the pool size, we don't put it back
	// バッファがプールサイズより小さい場合は戻さない
}

// Common packet buffer sizes
// 一般的なパケットバッファサイズ
const (
	SmallPacketSize  = 128   // For small headers / 小さいヘッダー用
	MediumPacketSize = 1500  // Typical MTU size / 一般的なMTUサイズ
	LargePacketSize  = 9000  // Jumbo frame size / ジャンボフレームサイズ
)

// Global byte slice pools for different sizes
// 異なるサイズのグローバルバイトスライスプール
var (
	smallBytesPool  = NewBytesPool(SmallPacketSize)
	mediumBytesPool = NewBytesPool(MediumPacketSize)
	largeBytesPool  = NewBytesPool(LargePacketSize)
)

// GetSmallBytes retrieves a small byte slice from the global pool
// グローバルプールから小さいバイトスライスを取得します
func GetSmallBytes() []byte {
	return smallBytesPool.Get()
}

// PutSmallBytes returns a small byte slice to the global pool
// 小さいバイトスライスをグローバルプールに返します
func PutSmallBytes(buf []byte) {
	smallBytesPool.Put(buf)
}

// GetMediumBytes retrieves a medium byte slice from the global pool
// グローバルプールから中サイズのバイトスライスを取得します
func GetMediumBytes() []byte {
	return mediumBytesPool.Get()
}

// PutMediumBytes returns a medium byte slice to the global pool
// 中サイズのバイトスライスをグローバルプールに返します
func PutMediumBytes(buf []byte) {
	mediumBytesPool.Put(buf)
}

// GetLargeBytes retrieves a large byte slice from the global pool
// グローバルプールから大きいバイトスライスを取得します
func GetLargeBytes() []byte {
	return largeBytesPool.Get()
}

// PutLargeBytes returns a large byte slice to the global pool
// 大きいバイトスライスをグローバルプールに返します
func PutLargeBytes(buf []byte) {
	largeBytesPool.Put(buf)
}

// GetBytes retrieves an appropriately sized byte slice based on the requested size
// 要求されたサイズに基づいて適切なサイズのバイトスライスを取得します
func GetBytes(size int) []byte {
	if size <= SmallPacketSize {
		return GetSmallBytes()
	} else if size <= MediumPacketSize {
		return GetMediumBytes()
	} else {
		return GetLargeBytes()
	}
}

// PutBytes returns a byte slice to the appropriate global pool based on its capacity
// バイトスライスをその容量に基づいて適切なグローバルプールに返します
func PutBytes(buf []byte) {
	capacity := cap(buf)
	if capacity <= SmallPacketSize {
		PutSmallBytes(buf)
	} else if capacity <= MediumPacketSize {
		PutMediumBytes(buf)
	} else if capacity <= LargePacketSize {
		PutLargeBytes(buf)
	}
	// If the buffer is larger than LargePacketSize, we don't put it back
	// バッファがLargePacketSizeより大きい場合は戻さない
}

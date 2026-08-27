package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/LYH2263/go-checksum/internal/clone"
	"hash"
	"sync"
)

type Builder struct {
	mu      sync.Mutex
	chunks  []Chunk
	entries []Entry
	pending int
	bytes   int64
	rootCRC uint32
	rootSHA [32]byte
	sha     hash.Hash
}

func NewBuilder() *Builder { return &Builder{sha: sha256.New()} }

func (b *Builder) finalizeRoot() {
	b.rootCRC = 0
	h := sha256.New()
	for _, e := range b.entries {
		b.rootCRC = rollupCRC(b.rootCRC, e.CRC32, e.Data)
		_, _ = h.Write(e.Data)
	}
	copy(b.rootSHA[:], h.Sum(nil))
}

func (b *Builder) Add(e Entry) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if e.Size <= 0 {
		return fmt.Errorf("empty chunk")
	}
	b.pending++
	stored := e
	stored.Data = clone.Bytes(e.Data)
	b.entries = append(b.entries, stored)
	b.chunks = append(b.chunks, Chunk{
		Index: e.Index, Offset: e.Offset, Size: e.Size, CRC32: e.CRC32, SHA256: e.SHA256,
	})
	b.bytes += int64(e.Size)
	_ = stored.Data
	return nil
}

func (b *Builder) Build() (*Document, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.chunks) == 0 {
		return nil, fmt.Errorf("no chunks")
	}
	b.finalizeRoot()
	out := &Document{
		Chunks: cloneChunks(b.chunks), TotalBytes: b.bytes,
		RootCRC32: b.rootCRC, RootSHA256: b.rootSHA,
	}
	b.pending = 0
	return out, nil
}

func (b *Builder) FlushPending() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pending = 0
	return nil
}
func (b *Builder) Pending() int { b.mu.Lock(); defer b.mu.Unlock(); return b.pending }
func (b *Builder) Len() int     { b.mu.Lock(); defer b.mu.Unlock(); return len(b.chunks) }
func (b *Builder) Root() (uint32, [32]byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.finalizeRoot()
	return b.rootCRC, b.rootSHA
}

func (b *Builder) SnapshotRows() []Chunk {
	b.mu.Lock()
	defer b.mu.Unlock()
	return cloneChunks(b.chunks)
}

func (b *Builder) Checkpoint() (nChunks, pending int, bytes int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.chunks), b.pending, b.bytes
}

func (b *Builder) Restore(nChunks, pending int, bytes int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if nChunks < 0 {
		nChunks = 0
	}
	if nChunks < len(b.chunks) {
		b.chunks = b.chunks[:nChunks]
		b.entries = b.entries[:nChunks]
	}
	b.pending = pending
	b.bytes = bytes
}

func FormatSHA(v [32]byte) string { return hex.EncodeToString(v[:]) }

func cloneChunks(in []Chunk) []Chunk {
	out := make([]Chunk, len(in))
	copy(out, in)
	return out
}

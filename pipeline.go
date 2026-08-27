package checksum

import (
	"context"
	"sync"

	"github.com/LYH2263/go-checksum/internal/chunk"
	"github.com/LYH2263/go-checksum/internal/manifest"
	"github.com/LYH2263/go-checksum/internal/rolling"
)

type Pipeline struct {
	mu       sync.Mutex
	opts     Options
	closed   bool
	splitter *chunk.Splitter
	builder  *manifest.Builder
	rolling  *rolling.Window
	total    int64
}

func New(opts Options) *Pipeline {
	opts = opts.withDefaults()
	return &Pipeline{
		opts: opts, splitter: chunk.NewSplitter(opts.ChunkSize),
		builder: manifest.NewBuilder(), rolling: rolling.NewWindow(),
	}
}

func (p *Pipeline) Stats() Stats {
	p.mu.Lock()
	defer p.mu.Unlock()
	rootCRC, rootSHA := p.builder.Root()
	return Stats{
		Chunks: p.builder.Len(), Bytes: p.total, Pending: p.builder.Pending(),
		RootCRC32: rootCRC, RootSHA256: manifest.FormatSHA(rootSHA),
	}
}

func (p *Pipeline) PendingChunks() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.builder.Pending()
}

func (p *Pipeline) Feed(ctx context.Context, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return ErrClosed
	}
	return p.feedLocked(ctx, data)
}

func (p *Pipeline) feedLocked(ctx context.Context, data []byte) error {
	off := p.total
	parts, err := p.splitter.Feed(ctx, data)
	if err != nil {
		return err
	}
	snapChunks, snapPending, snapBytes := p.builder.Checkpoint()
	snapCRC, snapN := p.rolling.Snapshot()
	for _, part := range parts {
		if err := ctx.Err(); err != nil {
			p.builder.Restore(snapChunks, snapPending, snapBytes)
			p.rolling.Restore(snapCRC, snapN)
			return err
		}
		h := NewMultiHasher()
		_, _ = h.Write(part.Data)
		crc, sha := h.Sum()
		if err := p.builder.Add(manifest.Entry{
			Index: part.Index, Offset: off + part.Offset, Size: len(part.Data),
			CRC32: crc, SHA256: sha, Data: append([]byte(nil), part.Data...),
		}); err != nil {
			p.builder.Restore(snapChunks, snapPending, snapBytes)
			p.rolling.Restore(snapCRC, snapN)
			return err
		}
		p.rolling.Update(part.Data)
	}
	p.total += int64(len(data))
	return nil
}

func (p *Pipeline) Finish() (*Manifest, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, ErrClosed
	}
	tail, err := p.splitter.Flush(ctxBackground())
	if err != nil {
		return nil, err
	}
	if tail != nil {
		h := NewMultiHasher()
		_, _ = h.Write(tail.Data)
		crc, sha := h.Sum()
		p.rolling.Update(tail.Data)
		if err := p.builder.Add(manifest.Entry{
			Index: tail.Index, Offset: p.total - int64(len(tail.Data)), Size: len(tail.Data),
			CRC32: crc, SHA256: sha, Data: append([]byte(nil), tail.Data...),
		}); err != nil {
			return nil, err
		}
	}
	return p.builder.Build()
}

func (p *Pipeline) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	if p.builder.Pending() > 0 {
		return ErrPending
	}
	p.closed = true
	return nil
}

func (p *Pipeline) CloseFlushCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return 0
	}
	// Capture the pending count BEFORE flushing: FlushPending zeroes the
	// counter, so reading it afterwards always returned 0 — the count fell
	// out of sync with the chunks still held in the builder's buffer.
	n := p.builder.Pending()
	_ = p.builder.FlushPending()
	p.closed = true
	return n
}

func (p *Pipeline) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.splitter.Reset()
	p.builder = manifest.NewBuilder()
	p.rolling.Reset()
	p.total = 0
	p.closed = false
}

func (p *Pipeline) RollingCRC() uint32 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rolling.CRC()
}

func ctxBackground() context.Context { return context.Background() }

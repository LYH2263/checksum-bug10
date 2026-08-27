package checksum

import (
	"bytes"
	"context"
	"fmt"

	"github.com/LYH2263/go-checksum/internal/manifest"
)

func (p *Pipeline) Verify(ctx context.Context, data []byte, doc *Manifest) error {
	if doc == nil {
		return fmt.Errorf("%w: nil manifest", ErrInvalid)
	}
	other, err := p.buildManifest(ctx, data)
	if err != nil {
		return err
	}
	if other.TotalBytes != doc.TotalBytes {
		return fmt.Errorf("%w: size %d vs %d", ErrMismatch, other.TotalBytes, doc.TotalBytes)
	}
	if other.RootCRC32 != doc.RootCRC32 {
		return fmt.Errorf("%w: root crc", ErrMismatch)
	}
	if !bytes.Equal(other.RootSHA256[:], doc.RootSHA256[:]) {
		return fmt.Errorf("%w: root sha", ErrMismatch)
	}
	if len(other.Chunks) != len(doc.Chunks) {
		return fmt.Errorf("%w: chunk count", ErrMismatch)
	}
	for i := range doc.Chunks {
		a, b := doc.Chunks[i], other.Chunks[i]
		if a.Index != b.Index || a.Size != b.Size || a.Offset != b.Offset {
			return fmt.Errorf("%w: chunk meta %d", ErrMismatch, i)
		}
		if a.CRC32 != b.CRC32 {
			return fmt.Errorf("%w: chunk crc %d", ErrMismatch, i)
		}
		if !bytes.Equal(a.SHA256[:], b.SHA256[:]) {
			return fmt.Errorf("%w: chunk sha %d", ErrMismatch, i)
		}
	}
	return nil
}

func (p *Pipeline) buildManifest(ctx context.Context, data []byte) (*manifest.Document, error) {
	q := New(p.opts)
	if err := q.Feed(ctx, data); err != nil {
		return nil, err
	}
	return q.Finish()
}

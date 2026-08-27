package checksum

import (
	"github.com/LYH2263/go-checksum/internal/codec"
	"github.com/LYH2263/go-checksum/internal/manifest"
)

type Manifest = manifest.Document

func EncodeManifest(m *Manifest) ([]byte, error) {
	return codec.EncodeJSON(m)
}

func DecodeManifest(b []byte) (*Manifest, error) {
	return codec.DecodeJSON(b)
}

func (p *Pipeline) SnapshotChunks() []ChunkView {
	p.mu.Lock()
	defer p.mu.Unlock()
	rows := p.builder.SnapshotRows()
	out := make([]ChunkView, len(rows))
	for i, c := range rows {
		out[i] = ChunkView{
			Index: c.Index, Offset: c.Offset, Size: c.Size,
			CRC32: c.CRC32, SHA256: manifest.FormatSHA(c.SHA256),
		}
	}
	return out
}

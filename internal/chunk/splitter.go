package chunk

import (
	"context"
	"fmt"
	"github.com/LYH2263/go-checksum/internal/clone"
)

type Part struct {
	Index  int
	Offset int64
	Data   []byte
}

type Splitter struct {
	size   int
	buf    []byte
	index  int
	offset int64
}

func NewSplitter(size int) *Splitter {
	if size <= 0 {
		size = 64 * 1024
	}
	return &Splitter{size: size}
}

func (s *Splitter) Feed(ctx context.Context, data []byte) ([]Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	startIdx, startOff, startBuf := s.index, s.offset, len(s.buf)
	remaining := data
	var out []Part
	for len(remaining) >= s.size {
		if err := ctx.Err(); err != nil {
			s.index, s.offset = startIdx, startOff
			if startBuf < len(s.buf) {
				s.buf = s.buf[:startBuf]
			}
			return nil, err
		}
		out = append(out, Part{
			Index: s.index, Offset: s.offset + int64(len(data)-len(remaining)),
			Data: clone.Bytes(remaining[:s.size]),
		})
		s.index++
		s.offset += int64(s.size)
		remaining = remaining[s.size:]
	}
	if len(remaining) > 0 {
		s.buf = append(s.buf, remaining...)
	}
	return out, nil
}

func (s *Splitter) Flush(ctx context.Context) (*Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(s.buf) == 0 {
		return nil, nil
	}
	p := &Part{Index: s.index, Offset: s.offset, Data: clone.Bytes(s.buf)}
	s.buf = nil
	s.index++
	return p, nil
}

func (s *Splitter) Reset()         { s.buf = nil; s.index = 0; s.offset = 0 }
func (s *Splitter) ChunkSize() int { return s.size }

func ValidateSize(n int) error {
	if n <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}
	return nil
}

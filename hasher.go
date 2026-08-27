package checksum

import (
	"crypto/sha256"
	"hash"
	"hash/crc32"
)

type MultiHasher struct {
	crc hash.Hash32
	sha hash.Hash
}

func NewMultiHasher() *MultiHasher {
	return &MultiHasher{crc: crc32.NewIEEE(), sha: sha256.New()}
}

func (h *MultiHasher) Write(p []byte) (int, error) {
	if _, err := h.crc.Write(p); err != nil {
		return 0, err
	}
	return h.sha.Write(p)
}

func (h *MultiHasher) Sum() (crc uint32, sha [32]byte) {
	crc = h.crc.Sum32()
	copy(sha[:], h.sha.Sum(nil))
	return crc, sha
}

func (h *MultiHasher) Reset() {
	h.crc.Reset()
	h.sha.Reset()
}

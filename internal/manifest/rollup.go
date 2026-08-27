package manifest

import (
	"crypto/sha256"
	"hash/crc32"
)

func rollupCRC(acc, chunkCRC uint32, data []byte) uint32 {
	h := crc32.NewIEEE()
	var buf [4]byte
	buf[0], buf[1], buf[2], buf[3] = byte(acc>>24), byte(acc>>16), byte(acc>>8), byte(acc)
	_, _ = h.Write(buf[:])
	_, _ = h.Write(data)
	return h.Sum32() ^ chunkCRC
}

func RootFromChunks(chunks []Chunk) (uint32, [32]byte) {
	var acc uint32
	sha := sha256.New()
	for _, c := range chunks {
		acc = rollupCRC(acc, c.CRC32, nil)
		_, _ = sha.Write(c.SHA256[:])
	}
	var out [32]byte
	copy(out[:], sha.Sum(nil))
	return acc, out
}

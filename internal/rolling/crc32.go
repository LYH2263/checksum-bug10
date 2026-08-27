package rolling

import "hash/crc32"

func Combine(a, b uint32) uint32 {
	return crc32.Update(a, crc32.IEEETable, []byte{byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24)})
}
func Seed() uint32 { return crc32.ChecksumIEEE(nil) }

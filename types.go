package checksum

type ChunkView struct {
	Index  int    `json:"index"`
	Offset int64  `json:"offset"`
	Size   int    `json:"size"`
	CRC32  uint32 `json:"crc32"`
	SHA256 string `json:"sha256"`
}

type Stats struct {
	Chunks     int    `json:"chunks"`
	Bytes      int64  `json:"bytes"`
	Pending    int    `json:"pending"`
	RootCRC32  uint32 `json:"root_crc32"`
	RootSHA256 string `json:"root_sha256"`
}

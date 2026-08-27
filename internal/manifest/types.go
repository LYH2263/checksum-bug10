package manifest

type Entry struct {
	Index  int
	Offset int64
	Size   int
	CRC32  uint32
	SHA256 [32]byte
	Data   []byte
}

type Chunk struct {
	Index  int      `json:"index"`
	Offset int64    `json:"offset"`
	Size   int      `json:"size"`
	CRC32  uint32   `json:"crc32"`
	SHA256 [32]byte `json:"sha256"`
}

type Document struct {
	Chunks     []Chunk  `json:"chunks"`
	TotalBytes int64    `json:"total_bytes"`
	RootCRC32  uint32   `json:"root_crc32"`
	RootSHA256 [32]byte `json:"root_sha256"`
}

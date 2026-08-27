package checksum

type Options struct {
	ChunkSize int
}

func (o Options) withDefaults() Options {
	if o.ChunkSize <= 0 {
		o.ChunkSize = 64 * 1024
	}
	return o
}

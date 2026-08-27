package rolling

import "hash/crc32"

type Window struct {
	crc uint32
	n   int64
}

func NewWindow() *Window { return &Window{} }
func (w *Window) Update(p []byte) {
	if len(p) == 0 {
		return
	}
	w.crc = crc32.Update(w.crc, crc32.IEEETable, p)
	w.n += int64(len(p))
}
func (w *Window) CRC() uint32                 { return w.crc }
func (w *Window) Bytes() int64                { return w.n }
func (w *Window) Reset()                      { w.crc = 0; w.n = 0 }
func (w *Window) Snapshot() (uint32, int64)   { return w.crc, w.n }
func (w *Window) Restore(crc uint32, n int64) { w.crc = crc; w.n = n }

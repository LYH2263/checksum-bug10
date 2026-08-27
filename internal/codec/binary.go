package codec

import (
	"bytes"
	"encoding/binary"
	"github.com/LYH2263/go-checksum/internal/manifest"
	"io"
)

func WriteBinary(w io.Writer, doc *manifest.Document) error {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, doc.TotalBytes); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(doc.Chunks))); err != nil {
		return err
	}
	for _, c := range doc.Chunks {
		if err := binary.Write(&buf, binary.LittleEndian, int32(c.Index)); err != nil {
			return err
		}
		if err := binary.Write(&buf, binary.LittleEndian, c.Offset); err != nil {
			return err
		}
		if err := binary.Write(&buf, binary.LittleEndian, int32(c.Size)); err != nil {
			return err
		}
		if err := binary.Write(&buf, binary.LittleEndian, c.CRC32); err != nil {
			return err
		}
		if _, err := buf.Write(c.SHA256[:]); err != nil {
			return err
		}
	}
	if err := binary.Write(&buf, binary.LittleEndian, doc.RootCRC32); err != nil {
		return err
	}
	if _, err := buf.Write(doc.RootSHA256[:]); err != nil {
		return err
	}
	_, err := w.Write(buf.Bytes())
	return err
}

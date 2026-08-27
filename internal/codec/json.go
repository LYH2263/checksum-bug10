package codec

import (
	"encoding/json"
	"github.com/LYH2263/go-checksum/internal/manifest"
	"io"
)

func EncodeJSON(doc *manifest.Document) ([]byte, error) { return json.Marshal(doc) }

func DecodeJSON(b []byte) (*manifest.Document, error) {
	var doc manifest.Document
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func DecodeJSONFrom(r io.Reader) (*manifest.Document, error) {
	dec := json.NewDecoder(r)
	var doc manifest.Document
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	if c, ok := r.(io.Closer); ok {
		defer c.Close()
	}
	return &doc, nil
}

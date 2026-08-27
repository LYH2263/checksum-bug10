package checksum

import (
	"context"
	"testing"
)

func TestBug10_CloseFlushCountOrder(t *testing.T) {
	p := New(Options{ChunkSize: 8})
	if err := p.Feed(context.Background(), []byte("12345678")); err != nil {
		t.Fatal(err)
	}
	n := p.CloseFlushCount()
	if n <= 0 {
		t.Fatalf("expected pending flush count >0, got %d", n)
	}
}

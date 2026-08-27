package main

import (
	"flag"
	"github.com/LYH2263/go-checksum"
	"github.com/LYH2263/go-checksum/inspector"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":8225", "listen address")
	chunk := flag.Int("chunk", 65536, "chunk size")
	flag.Parse()
	p := checksum.New(checksum.Options{ChunkSize: *chunk})
	api := &inspector.API{Pipe: p}
	srv := inspector.New(api)
	log.Printf("checksum inspector on %s", *addr)
	if err := http.ListenAndServe(*addr, srv.Handler); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

package inspector

import (
	"context"
	"encoding/json"
	"github.com/LYH2263/go-checksum"
	"io"
	"net/http"
)

type API struct{ Pipe *checksum.Pipeline }

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.URL.Path {
	case "/api/stats":
		writeJSON(w, a.Pipe.Stats())
	case "/api/chunks":
		writeJSON(w, a.Pipe.SnapshotChunks())
	case "/api/digest":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		doc, err := checksum.ServeDigest(ctx, a.Pipe, body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, doc)
	default:
		http.NotFound(w, r)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

var _ context.Context

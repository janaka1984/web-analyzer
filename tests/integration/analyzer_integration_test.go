package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/janaka/web-analyzer/internal/adapter/httpfetch"
	"github.com/janaka/web-analyzer/internal/service/analyzer"
)

func TestFullAnalyzeAgainstTestServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!doctype html><html><head><title>Home</title></head>
		<body><h1>Welcome</h1><a href="/ok">ok</a><a href="/nope">broken</a></body></html>`))
	})
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	// /nope -> no handler => 404

	ts := httptest.NewServer(mux)
	defer ts.Close()

	fetch := httpfetch.New(3, 3, 5)
	ana := analyzer.NewDefaultAnalyzer(fetch, 5, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := ana.Analyze(ctx, ts.URL+"/index.html")
	if err != nil {
		t.Fatalf("analyze err: %v", err)
	}
	if res.Title != "Home" {
		t.Fatalf("title: %s", res.Title)
	}
	if res.H1 != 1 {
		t.Fatalf("h1: %d", res.H1)
	}
	if res.InternalLinks != 2 {
		t.Fatalf("internal: %d", res.InternalLinks)
	}
	if res.BrokenLinks < 1 {
		t.Fatalf("broken links expected >= 1, got %d", res.BrokenLinks)
	}
}

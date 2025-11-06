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
		if _, err := w.Write([]byte(`<!doctype html><html><head><title>Home</title></head>`)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}

	})
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

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

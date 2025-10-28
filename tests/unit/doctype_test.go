package unit

import (
	"testing"

	"github.com/janaka/web-analyzer/internal/service/analyzer"
)

func TestDetectHTMLVersion(t *testing.T) {
	html5 := []byte("<!doctype html><html><head><title>x</title></head><body></body></html>")
	got := analyzer.DetectHTMLVersion(html5)
	if got != "HTML5" {
		t.Fatalf("want HTML5, got %s", got)
	}
}

package unit

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/janaka/web-analyzer/internal/service/analyzer"
)

func TestExtractLinks(t *testing.T) {
	html := `<a href="/a">A</a><a href="https://ext.com/b">B</a>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	links := analyzer.ExtractLinks(doc, "https://example.com/base")
	if len(links) != 2 {
		t.Fatalf("want 2 links, got %d", len(links))
	}
	if links[0] != "https://example.com/a" {
		t.Fatalf("resolved: %s", links[0])
	}
}

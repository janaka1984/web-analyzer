package unit

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/janaka/web-analyzer/internal/service/analyzer"
)

func TestHeadings(t *testing.T) {
	s := `<html><body><h1>a</h1><h2>b</h2><h2>c</h2><h3>d</h3></body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(s))
	h := analyzer.CountHeadings(doc)
	if h[1] != 1 || h[2] != 2 || h[3] != 1 {
		t.Fatalf("unexpected headings: %#v", h)
	}
}

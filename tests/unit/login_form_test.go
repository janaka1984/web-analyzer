package unit

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/janaka/web-analyzer/internal/service/analyzer"
)

func TestHasLoginForm(t *testing.T) {
	html := `<form><input type="password"/></form>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	if !analyzer.HasLoginForm(doc) {
		t.Fatal("expected login form detection")
	}
}

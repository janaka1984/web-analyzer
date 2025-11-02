package analyzer

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/janaka/web-analyzer/pkg/validator"
)

func ExtractLinks(doc *goquery.Document, base string) []string {
	var out []string
	seen := make(map[string]bool)

	b, err := url.Parse(base)
	if err != nil {
		return out
	}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		resolved := b.ResolveReference(u).String()

		// Normalize *after* resolution
		norm := validator.NormalizeLink(resolved)

		// De-duplicate
		if seen[norm] {
			return
		}
		seen[norm] = true

		out = append(out, norm)
	})

	return out
}

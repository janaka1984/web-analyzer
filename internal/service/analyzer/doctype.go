package analyzer

import (
	"bytes"
	"regexp"
)

var reDoctype = regexp.MustCompile(`(?i)<!DOCTYPE\s+html(?:\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+([0-9.]+)[^"]*|[^>]*)>`)

// DetectHTMLVersion determines HTML version (HTML5, HTML4.x, XHTML, or Unknown).
func DetectHTMLVersion(body []byte) string {
	s := bytes.ToLower(bytes.TrimSpace(body))

	// Explicit HTML5 check first
	if bytes.Contains(s, []byte("<!doctype html>")) {
		return "HTML5"
	}

	// Check for HTML 4.x / XHTML
	if m := reDoctype.FindSubmatch(body); len(m) >= 2 {
		return "HTML " + string(m[1])
	}

	// Generic XHTML check
	if bytes.Contains(s, []byte("xhtml")) {
		return "XHTML"
	}

	// Fallback
	if bytes.Contains(s, []byte("<!doctype")) {
		return "HTML (doctype present)"
	}

	return "Unknown/No DOCTYPE"
}

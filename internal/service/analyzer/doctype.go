package analyzer

import (
	"bytes"
	"regexp"
)

var reDoctype = regexp.MustCompile(`(?i)<!DOCTYPE\s+html(?:\s+PUBLIC\s+"-//W3C//DTD\s+HTML\s+([0-9.]+)[^"]*|[^>]*)>`)

// DetectHTMLVersion uses a simple heuristic on DOCTYPE / legacy markers.
func DetectHTMLVersion(body []byte) string {
	s := bytes.ToUpper(bytes.TrimSpace(body))
	switch {
	case bytes.Contains(s, []byte("<!DOCTYPE HTML>")), bytes.Contains(s, []byte("<!DOCTYPE HTML PUBLIC")):
		// Try to extract version from XHTML/HTML4 doctypes
		if m := reDoctype.FindSubmatch(body); len(m) >= 2 {
			return "HTML " + string(m[1])
		}
		// HTML5 doctype is <!DOCTYPE html>
		if bytes.Contains(bytes.ToLower(body), []byte("<!doctype html>")) {
			return "HTML5"
		}
		return "HTML (doctype present)"
	case bytes.Contains(bytes.ToLower(s), []byte("<!doctype html>")):
		return "HTML5"
	default:
		return "Unknown/No DOCTYPE"
	}
}

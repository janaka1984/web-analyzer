package analyzer

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var loginTextRe = regexp.MustCompile(`(?i)\b(log\s?in|sign\s?in|connexion|anmelden|accedi)\b`)

func HasLoginForm(doc *goquery.Document) bool {
	// 1) Classic password input
	if doc.Find(`input[type="password"], input[type=password]`).Length() > 0 {
		return true
	}

	// 2) Heuristic matches
	if doc.Find(`input[name*="pass"], input[id*="pass"], input[name*="pwd"], input[id*="pwd"]`).Length() > 0 {
		return true
	}

	// 3) Forms with login-related IDs/actions
	if doc.Find(`form[action*="login"], form[id*="login"], form[class*="login"]`).Length() > 0 {
		return true
	}

	// 4) Visible Login buttons/links
	found := false
	doc.Find("a,button,input[type='submit']").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := strings.ToLower(strings.TrimSpace(s.Text()))
		val, _ := s.Attr("value")
		href, _ := s.Attr("href")
		if loginTextRe.MatchString(text) || loginTextRe.MatchString(val) || loginTextRe.MatchString(href) {
			found = true
			return false
		}
		return true
	})
	return found
}

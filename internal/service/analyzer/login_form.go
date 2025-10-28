package analyzer

import "github.com/PuerkitoBio/goquery"

// Heuristic: a login form typically contains an <input type="password">.
func HasLoginForm(doc *goquery.Document) bool {
	return doc.Find("form input[type='password']").Length() > 0
}

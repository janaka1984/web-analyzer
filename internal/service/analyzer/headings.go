package analyzer

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func CountHeadings(doc *goquery.Document) map[int]int {
	out := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0}
	for i := 1; i <= 6; i++ {
		out[i] = doc.Find("h" + strconv.Itoa(i)).Length()
	}
	return out
}

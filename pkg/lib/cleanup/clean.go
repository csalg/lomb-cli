package cleanup

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RemoveCSS removes blacklisted CSS from the document.
func (a *ParsedHTML) RemoveCSS(blacklist []string) {
	doc := a.GoQueryDoc()
	if blacklist == nil {
		return
	}
	blacklistStr := strings.Join(blacklist, ",")
	// find first blacklisted match, remove it from the document
	for {
		blackSel := doc.FindMatcher(goquery.Single(blacklistStr))
		if blackSel == nil || len(blackSel.Nodes) == 0 {
			break
		}
		blackSel.Remove()
	}
}

package study

import "github.com/csalg/lomb-cli/pkg/types"

type PageModel struct {
	View           types.View
	Data           Data
	AvailableViews []string
}

type Data struct {
	LemmaCounts      []LemmaCount
	ReaderParagraphs []types.Paragraph
}

type LemmaCount struct {
	Lemma string
	Count int
}

type Example struct {
	Text        string
	Translation string
}

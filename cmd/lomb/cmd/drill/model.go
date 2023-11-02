package drill

import "github.com/csalg/lomb-cli/pkg/types"

type PageModel struct {
	Grid [][]Cell
	Data Data
}

type Cell struct {
	Type          CellType
	ID            string
	DictionaryURL string
}

type CellType string

const (
	CellLemmaCounts             CellType = "lemma-counts"
	CellExamples                CellType = "examples"
	CellDictionary              CellType = "dictionary"
	CellUnderstandableSentences CellType = "understandable-sentences"
	CellReader                  CellType = "reader"
)

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

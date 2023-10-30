package drill

type PageModel struct {
	Grid Grid
	Data Data
}

type Grid struct {
	Columns []Column
}

type Column struct {
	Rows []Row
}

type Row struct {
	Type          RowType
	ID            string
	DictionaryURL string
}

type RowType string

const (
	RowTypeLemmaCounts RowType = "lemma-counts"
	RowTypeExamples    RowType = "examples"
	RowTypeDictionary  RowType = "dictionary"
)

type Data struct {
	LemmaCounts []LemmaCount
}

type LemmaCount struct {
	Lemma string
	Count int
}

type Example struct {
	Text        string
	Translation string
}

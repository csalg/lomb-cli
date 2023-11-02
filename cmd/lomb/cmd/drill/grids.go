package drill

var (
	DefaultGrid = [][]Cell{
		// Column 1
		{
			{
				Type: CellLemmaCounts,
				ID:   "lemma-counts",
			},
			{
				Type: CellExamples,
				ID:   "examples",
			},
		},
		// Column 2
		{
			{
				Type:          CellDictionary,
				ID:            "dictcom",
				DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
			},
		},
		// Column 3
		{
			{
				Type:          CellDictionary,
				ID:            "wiktionary",
				DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
			},
		},
	}
	UnderstandableSentences = [][]Cell{
		// Column 1
		{
			{
				Type: CellUnderstandableSentences,
				ID:   "understandable-sentences",
			},
		},
	}
)

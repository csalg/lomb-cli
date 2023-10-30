package drill

var (
	DefaultGrid = Grid{
		Columns: []Column{
			{
				Rows: []Row{
					{
						Type: RowTypeLemmaCounts,
						ID:   "lemma-counts",
					},
					{
						Type: RowTypeExamples,
						ID:   "examples",
					},
				},
			},
			{
				Rows: []Row{
					{
						Type:          RowTypeDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
				},
			},
			{
				Rows: []Row{
					{
						Type:          RowTypeDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
	}
	SimpleSentencesGrid = Grid{
		Columns: []Column{
			{
				Rows: []Row{
					{
						Type: RowTypeLemmaCounts,
						ID:   "lemma-counts",
					},
					{
						Type: RowTypeExamples,
						ID:   "examples",
					},
				},
			},
			{
				Rows: []Row{
					{
						Type:          RowTypeDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
				},
			},
			{
				Rows: []Row{
					{
						Type:          RowTypeDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
	}
)

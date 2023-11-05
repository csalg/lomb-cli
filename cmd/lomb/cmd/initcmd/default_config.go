package initcmd

import "github.com/csalg/lomb-cli/pkg/types"

var DefaultLanguagePairConfig = types.LanguagePairConfig{
	BaseLanguage:   "es",
	TargetLanguage: "bg",
	Views: []types.View{
		{
			Name: "lemma-counts-and-examples",
			Grid: [][]types.Cell{
				// Column 1
				{
					{
						Type: types.CellLemmaCounts,
						ID:   "lemma-counts",
					},
					{
						Type: types.CellExamples,
						ID:   "examples",
					},
				},
				// Column 2
				{
					{
						Type:          types.CellDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
				},
				// Column 3
				{
					{
						Type:          types.CellDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
		{
			Name: "all-data",
			Grid: [][]types.Cell{
				// Column 1
				{
					{
						Type: types.CellUnderstandableSentences,
						ID:   "understandable-sentences",
					},
				},
				// Column 2
				{
					{
						Type: types.CellLemmaCounts,
						ID:   "lemma-counts",
					},
					{
						Type: types.CellExamples,
						ID:   "examples",
					},
				},
				// Column 3
				{
					{
						Type: types.CellReader,
						ID:   "reader",
					},
				},
			},
		},
		{
			Name: "all-data-with-dictionaries",
			Grid: [][]types.Cell{
				// Column 1
				{

					{
						Type: types.CellReader,
						ID:   "examples",
					},
					{
						Type: types.CellUnderstandableSentences,
						ID:   "understandable-sentences",
					},
				},
				// Column 2
				{
					{
						Type: types.CellLemmaCounts,
						ID:   "lemma-counts",
					},
					{
						Type: types.CellExamples,
						ID:   "examples",
					},
				},
				// Column 3
				{
					{
						Type:          types.CellDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
					{
						Type:          types.CellDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
		{
			Name: "just-lemmas",
			Grid: [][]types.Cell{
				// Column 1
				{
					{
						Type: types.CellLemmaCounts,
						ID:   "lemma-counts",
					},
					{
						Type: types.CellExamples,
						ID:   "examples",
					},
				},
				// Column 2
				{
					{
						Type:          types.CellDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
					{
						Type:          types.CellDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
		{
			Name: "just-read",
			Grid: [][]types.Cell{
				// Column 1
				{
					{
						Type: types.CellReader,
						ID:   "lemma-counts",
					},
				},
				// Column 2
				{
					{
						Type:          types.CellDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
					{
						Type:          types.CellDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
		{
			Name: "just-understandable-sentences",
			Grid: [][]types.Cell{
				// Column 1
				{
					{
						Type: types.CellUnderstandableSentences,
						ID:   "understandable-sentences",
					},
				},
				// Column 2
				{
					{
						Type:          types.CellDictionary,
						ID:            "dictcom",
						DictionaryURL: "https://www.dict.com/?t=bg&set=_bgen&w=$LEMMA",
					},
					{
						Type:          types.CellDictionary,
						ID:            "wiktionary",
						DictionaryURL: "https://en.wiktionary.org/w/index.php?title=$LEMMA#Bulgarian",
					},
				},
			},
		},
	},
}

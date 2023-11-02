package drill

import (
	"regexp"
	"sort"
	"strings"

	"github.com/csalg/lomb-cli/pkg/types"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	numbersRegex         = regexp.MustCompile(`\d+`)
)

type Corpus struct {
	LemmaCounts                    []LemmaCount
	LemmaExamplesLookup            map[string][]string
	TranslationLookup              map[string]string
	SentencesWithUnderstandability []SentenceWithUnderstandability
}

type SentenceWithUnderstandability struct {
	Text              string
	Understandability int
	Length            int
}

func NewCorpus() *Corpus {
	return &Corpus{
		LemmaCounts:                    make([]LemmaCount, 0),
		LemmaExamplesLookup:            make(map[string][]string),
		TranslationLookup:              make(map[string]string),
		SentencesWithUnderstandability: make([]SentenceWithUnderstandability, 0),
	}
}

// LoadTexts loads the given texts into the corpus.
func (crp *Corpus) LoadTexts(txts ...types.ProcessedText) {
	for _, txt := range txts {
		for _, paragraph := range txt.Paragraphs {
			for _, chunk := range paragraph {
				// Add translation to lookup
				crp.TranslationLookup[chunk.Sentence()] = chunk.Translation

				// Add examples to lookup
				for _, token := range chunk.Tokens {
					// Trim alphanumeric characters from the lemma.
					sanitizedLemma := sanitizeLemma(token.Lemma)
					if sanitizedLemma == "" {
						continue
					}
					crp.LemmaExamplesLookup[sanitizedLemma] = append(crp.LemmaExamplesLookup[sanitizedLemma], chunk.Sentence())
				}
			}
		}
	}

	// Re-calculate lemma counts
	var lemmaCounts []LemmaCount
	for lemma, examples := range crp.LemmaExamplesLookup {
		lemmaCounts = append(lemmaCounts, LemmaCount{Lemma: lemma, Count: len(examples)})
	}
	sort.Slice(lemmaCounts, func(i, j int) bool {
		return lemmaCounts[i].Count > lemmaCounts[j].Count
	})
	crp.LemmaCounts = lemmaCounts

	// Re-calculate chunk understandability
	crp.SentencesWithUnderstandability = make([]SentenceWithUnderstandability, 0, len(crp.TranslationLookup))
	for _, paragraph := range txts[0].Paragraphs {
		for _, chunk := range paragraph {
			mostInfrequentToken := 1000000000000
			for _, token := range chunk.Tokens {
				frequency := crp.LemmaExamplesLookup[sanitizeLemma(token.Lemma)]
				if len(frequency) < mostInfrequentToken {
					mostInfrequentToken = len(frequency)
				}
			}
			crp.SentencesWithUnderstandability = append(crp.SentencesWithUnderstandability, SentenceWithUnderstandability{
				Text:              chunk.Sentence(),
				Length:            len(chunk.Tokens),
				Understandability: mostInfrequentToken,
			})
		}
	}
}

// Translate translates the given string.
func (crp *Corpus) Translate(str string) string {
	return crp.TranslationLookup[str]
}

// GetExamples finds examples for the given lemma.
func (crp *Corpus) GetExamples(lemma string) []string {
	examplesNoDuplicates := make([]string, 0, len(crp.LemmaExamplesLookup[lemma]))
	seen := make(map[string]struct{})
	for _, example := range crp.LemmaExamplesLookup[lemma] {
		if _, ok := seen[example]; ok {
			continue
		}
		seen[example] = struct{}{}
		examplesNoDuplicates = append(examplesNoDuplicates, example)
	}
	// Sort by length
	sort.Slice(examplesNoDuplicates, func(i, j int) bool {
		return len(examplesNoDuplicates[i]) < len(examplesNoDuplicates[j])
	})
	return examplesNoDuplicates
}

// GetUnderstandableSentences finds understandable sentences
func (crp *Corpus) GetUnderstandableSentences(maxLength, minUnderstandability int) []SentenceWithUnderstandability {
	sentencesDedup := make([]SentenceWithUnderstandability, 0, len(crp.SentencesWithUnderstandability))
	seen := make(map[string]struct{})
	for _, snt := range crp.SentencesWithUnderstandability {
		if _, ok := seen[snt.Text]; ok {
			continue
		}
		if snt.Length > maxLength {
			continue
		}
		if snt.Understandability < minUnderstandability {
			continue
		}
		seen[snt.Text] = struct{}{}
		sentencesDedup = append(sentencesDedup, snt)
	}
	// Sort by understandability
	sort.Slice(sentencesDedup, func(i, j int) bool {
		return sentencesDedup[i].Understandability < sentencesDedup[j].Understandability
	})
	return sentencesDedup
}

func sanitizeLemma(str string) string {
	cleanStr := nonAlphanumericRegex.ReplaceAllString(str, "")
	cleanStr = numbersRegex.ReplaceAllString(cleanStr, "")
	return strings.TrimSpace(strings.ToLower(cleanStr))
}

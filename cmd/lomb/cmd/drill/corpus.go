package drill

import (
	"fmt"
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
	Paragraphs                     []types.Paragraph
	Stats                          Stats
}

type SentenceWithUnderstandability struct {
	Text              string
	Understandability int
	Length            int
}

type Stats struct {
	LengthAverage                   float64
	LengthStdDev                    float64
	MinLemmaCountPerSentenceAverage float64
	MinLemmaCountPerSentenceStdDev  float64
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
			crp.Paragraphs = append(crp.Paragraphs, paragraph)
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
	// Calculate understandability stats
	crp.Stats = Stats{}
	lengthSum := 0
	minLemmaCountPerSentenceSum := 0
	for _, sentence := range crp.SentencesWithUnderstandability {
		lengthSum += sentence.Length
		minLemmaCountPerSentenceSum += sentence.Understandability
	}
	crp.Stats.LengthAverage = float64(lengthSum) / float64(len(crp.SentencesWithUnderstandability))
	crp.Stats.MinLemmaCountPerSentenceAverage = float64(minLemmaCountPerSentenceSum) / float64(len(crp.SentencesWithUnderstandability))
	lengthVarianceSum := 0.0
	minLemmaCountPerSentenceVarianceSum := 0.0
	for _, sentence := range crp.SentencesWithUnderstandability {
		lengthVarianceSum += (float64(sentence.Length) - crp.Stats.LengthAverage) * (float64(sentence.Length) - crp.Stats.LengthAverage)
		minLemmaCountPerSentenceVarianceSum += (float64(sentence.Understandability) - crp.Stats.MinLemmaCountPerSentenceAverage) * (float64(sentence.Understandability) - crp.Stats.MinLemmaCountPerSentenceAverage)
	}
	crp.Stats.LengthStdDev = lengthVarianceSum / float64(len(crp.SentencesWithUnderstandability))
	crp.Stats.MinLemmaCountPerSentenceStdDev = minLemmaCountPerSentenceVarianceSum / float64(len(crp.SentencesWithUnderstandability))
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
func (crp *Corpus) GetUnderstandableSentences(maxLengthStdDev, minUnderstandabilityStdDev float64) []SentenceWithUnderstandability {
	maxLength := int(crp.Stats.LengthAverage + maxLengthStdDev*crp.Stats.LengthStdDev)
	minUnderstandability := int(crp.Stats.MinLemmaCountPerSentenceAverage + minUnderstandabilityStdDev*crp.Stats.MinLemmaCountPerSentenceStdDev)
	fmt.Println("maxLength:", maxLength)
	fmt.Println("minUnderstandability:", minUnderstandability)
	fmt.Println("crp.Stats.LengthAverage:", crp.Stats.LengthAverage)
	fmt.Println("crp.Stats.LengthStdDev:", crp.Stats.LengthStdDev)
	fmt.Println("crp.Stats.MinLemmaCountPerSentenceAverage:", crp.Stats.MinLemmaCountPerSentenceAverage)
	fmt.Println("crp.Stats.MinLemmaCountPerSentenceStdDev:", crp.Stats.MinLemmaCountPerSentenceStdDev)
	fmt.Println("len(crp.SentencesWithUnderstandability):", len(crp.SentencesWithUnderstandability))

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

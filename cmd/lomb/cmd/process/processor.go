package process

import (
	"fmt"
	"strings"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process/lemmatizers/dummylemmatizer"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process/lemmatizers/golem"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process/translators"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process/translators/deepl"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process/translators/dummytranslator"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils/assert"
	"github.com/csalg/lomb-cli/pkg/utils/itertools"
)

type TextProcessor struct {
	Lemmatizer Lemmatizer
	Translator Translator
	Config     Config
}

type Lemmatizer interface {
	Lemmatize(text string) ([]types.Token, error)
}

type Translator interface {
	Translate(sourceLang, targetLang types.Language, text []string) ([]translators.TranslatedText, error)
}

type Config struct {
	BaseLanguage          types.Language
	TargetLanguage        types.Language
	DeeplAPIKey           string
	DeeplAPIPro           bool
	GoogleTranslateAPIKey string
}

func NewTextProcessor(conf Config) (TextProcessor, error) {
	tp := TextProcessor{Config: conf}

	var err error
	// Set lemmatizer
	switch {
	case golem.IsLanguageSupported(conf.TargetLanguage):
		tp.Lemmatizer, err = golem.New(conf.TargetLanguage)
		if err != nil {
			return tp, fmt.Errorf("creating golem lemmatizer: %w", err)
		}
	default:
		tp.Lemmatizer = dummylemmatizer.New()
	}

	// Set translator
	switch {
	case conf.BaseLanguage == conf.TargetLanguage:
		tp.Translator = dummytranslator.New()
	case deepl.IsLanguageSupported(conf.BaseLanguage) && deepl.IsLanguageSupported(conf.TargetLanguage):
		tp.Translator = deepl.New(conf.DeeplAPIKey, conf.DeeplAPIPro)
	default:
		tp.Translator = dummytranslator.New()
	}
	return tp, nil
}

func (tp TextProcessor) Process(text string) (types.ProcessedText, error) {
	assert.NotNil(tp.Lemmatizer, "lemmatizer")
	assert.NotNil(tp.Translator, "translator")

	processedText := types.ProcessedText{
		BaseLanguage:   tp.Config.BaseLanguage,
		TargetLanguage: tp.Config.TargetLanguage,
	}

	// Lemmatization
	for _, paragraphStr := range strings.Split(text, "\n") {
		paragraph := types.Paragraph{}
		for _, sentence := range strings.Split(paragraphStr, ".") {
			tokens, err := tp.Lemmatizer.Lemmatize(sentence)
			if err != nil {
				return types.ProcessedText{}, fmt.Errorf("lemmatizing sentence %s: %w", sentence, err)
			}
			chunk := types.Chunk{
				Tokens: tokens,
			}
			paragraph = append(paragraph, chunk)
		}
		processedText.Paragraphs = append(processedText.Paragraphs, paragraph)
	}

	// Translation
	sentences := []string{}
	for _, paragraph := range processedText.Paragraphs {
		for _, chunk := range paragraph {
			sentences = append(sentences, chunk.Sentence())
		}
	}
	translations := make(map[string]string)
	var err error
	itertools.Chunk(sentences, 100, func(sentencesChunk []string, i int) bool {
		tr, translationErr := tp.Translator.Translate(tp.Config.TargetLanguage, tp.Config.BaseLanguage, sentencesChunk)
		if translationErr != nil {
			err = fmt.Errorf("translating sentences %v: %w", sentencesChunk, err)
			return false
		}
		for _, t := range tr {
			translations[t.Original] = t.Translated
		}
		return true
	})
	if err != nil {
		return types.ProcessedText{}, err
	}
	for i, paragraph := range processedText.Paragraphs {
		for j, chunk := range paragraph {
			translatedSentence, ok := translations[chunk.Sentence()]
			if !ok {
				return types.ProcessedText{}, fmt.Errorf("sentence %s not found in translations", chunk.Sentence())
			}
			processedText.Paragraphs[i][j].Translation = translatedSentence
		}
	}

	return processedText, nil
}

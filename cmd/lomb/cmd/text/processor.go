package text

import (
	"fmt"
	"strings"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/lemmatizers/dummylemmatizer"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/lemmatizers/golem"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/translators"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/translators/dummytranslator"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/translators/gpttranslator"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils/assert"
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
	OpenAIAPIKey          string
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
		// DeepL deprecated: Too expensive!
	// case deepl.IsLanguageSupported(conf.BaseLanguage) && deepl.IsLanguageSupported(conf.TargetLanguage):
	// 	tp.Translator = deepl.New(conf.DeeplAPIKey, conf.DeeplAPIPro)
	default:
		tp.Translator = gpttranslator.New(conf.OpenAIAPIKey)
		// tp.Translator = dummytranslator.New()
	}
	return tp, nil
}

type ProcessOptions struct {
	FromTranslation bool
}

func (tp TextProcessor) Process(text string, opts ProcessOptions) (types.ProcessedText, error) {
	assert.NotNil(tp.Lemmatizer, "lemmatizer")
	assert.NotNil(tp.Translator, "translator")

	processedText := types.ProcessedText{
		BaseLanguage:   tp.Config.BaseLanguage,
		TargetLanguage: tp.Config.TargetLanguage,
	}

	// Split text into paragraphs and chunks
	for _, paragraphStr := range strings.SplitAfter(text, "\n") {
		paragraph := types.Paragraph{}
		for _, sentence := range strings.SplitAfter(paragraphStr, ".") {
			chunk := types.Chunk{}
			if opts.FromTranslation {
				chunk.Translation = sentence
			} else {
				chunk.Text = sentence
			}
			paragraph = append(paragraph, chunk)
		}
		processedText.Paragraphs = append(processedText.Paragraphs, paragraph)
	}

	// Translation
	for i := range processedText.Paragraphs {
		for j := range processedText.Paragraphs[i] {
			if opts.FromTranslation {
				text, err := tp.Translator.Translate(tp.Config.BaseLanguage, tp.Config.TargetLanguage, []string{processedText.Paragraphs[i][j].Translation})
				if err != nil {
					return types.ProcessedText{}, fmt.Errorf("translating sentence %s: %w", processedText.Paragraphs[i][j].Translation, err)
				}
				processedText.Paragraphs[i][j].Text = text[0].Translated
			} else {
				text, err := tp.Translator.Translate(tp.Config.TargetLanguage, tp.Config.BaseLanguage, []string{processedText.Paragraphs[i][j].Text})
				if err != nil {
					return types.ProcessedText{}, fmt.Errorf("translating sentence %s: %w", processedText.Paragraphs[i][j].Text, err)
				}
				processedText.Paragraphs[i][j].Translation = text[0].Translated

			}
		}
	}

	// Lemmatization
	for i := range processedText.Paragraphs {
		for j := range processedText.Paragraphs[i] {
			tokens, err := tp.Lemmatizer.Lemmatize(processedText.Paragraphs[i][j].Text)
			if err != nil {
				return types.ProcessedText{}, fmt.Errorf("lemmatizing sentence %s: %w", processedText.Paragraphs[i][j].Text, err)
			}
			processedText.Paragraphs[i][j].Tokens = tokens
		}
	}
	return processedText, nil
}

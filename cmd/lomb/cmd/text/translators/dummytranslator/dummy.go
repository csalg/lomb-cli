package dummytranslator

import (
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/translators"
)

type DummyTranslator struct{}

func New() DummyTranslator {
	return DummyTranslator{}
}

func (dt DummyTranslator) Translate(sourceLang, targetLang string, text []string) ([]translators.TranslatedText, error) {
	var translatedTexts []translators.TranslatedText
	for _, t := range text {
		translatedTexts = append(translatedTexts, translators.TranslatedText{
			Translated: t,
			Original:   t,
		})
	}
	return translatedTexts, nil
}

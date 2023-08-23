package dummylemmatizer

import (
	"strings"

	"github.com/csalg/lomb-cli/pkg/types"
)

type DummyLemmatizer struct{}

func New() DummyLemmatizer {
	return DummyLemmatizer{}
}

func (dl DummyLemmatizer) IsSupported(lang types.Language) bool {
	return true
}

func (dl DummyLemmatizer) Lemmatize(text string) ([]types.Token, error) {
	tokens := []types.Token{}
	for _, word := range strings.Split(text, " ") {
		tokens = append(tokens, types.Token{
			Text:  word,
			Lemma: word,
		})
	}
	return tokens, nil
}

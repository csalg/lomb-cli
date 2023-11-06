package golem

import (
	"fmt"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/csalg/lomb-cli/pkg/types"
)

func IsLanguageSupported(language string) bool {
	return false
}

type GolemLemmatizer struct {
	lemmatizer *golem.Lemmatizer
}

func New(language string) (*GolemLemmatizer, error) {
	if !IsLanguageSupported(language) {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	gl := GolemLemmatizer{}
	return &gl, nil
}

func (gl GolemLemmatizer) Lemmatize(text string) ([]types.Token, error) {
	tokens := []types.Token{}
	for _, word := range strings.Split(text, " ") {
		cleanWord := strings.Trim(word, ".,;:!?\"'()[]{}")
		tokens = append(tokens, types.Token{
			Text:  word,
			Lemma: gl.lemmatizer.Lemma(strings.ToLower(cleanWord)),
		})
	}
	return tokens, nil
}

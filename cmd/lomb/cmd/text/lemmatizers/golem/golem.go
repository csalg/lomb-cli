package golem

import (
	"fmt"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/csalg/lomb-cli/pkg/types"
)

func IsLanguageSupported(language types.Language) bool {
	switch language {
	case types.English:
		return true
	default:
		return false
	}
}

type GolemLemmatizer struct {
	lemmatizer *golem.Lemmatizer
}

func New(language types.Language) (*GolemLemmatizer, error) {
	if !IsLanguageSupported(language) {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	gl := GolemLemmatizer{}
	var err error
	switch language {
	case types.English:
		gl.lemmatizer, err = golem.New(en.New())
		if err != nil {
			return nil, fmt.Errorf("creating golem lemmatizer: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
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

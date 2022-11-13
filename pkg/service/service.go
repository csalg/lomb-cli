package service

import (
	"fmt"
	"strings"

	"github.com/csalg/lomb-cli/pkg/types"
)

type Service struct {
	config types.Config
}

func New(config types.Config) Service {
	return Service{
		config: config,
	}
}

type GetDictionaryURLInput struct {
	SourceLanguage string
	TargetLanguage string
	Lemma          string
}

// GetDictionaryURLs returns the dictionary url for a given source language, target language, lemma
// triplet.
func (s Service) GetDictionaryURL(input *GetDictionaryURLInput) (string, error) {
	sourceLangApiName := ""
	targetLangApiName := ""

	for _, dictionary := range s.config.Dictionaries {
		for _, sourceLanguage := range dictionary.SourceLanguages {
			if sourceLanguage.Language == input.SourceLanguage {
				sourceLangApiName = sourceLanguage.APIName
			}
		}
		for _, targetLanguage := range dictionary.TargetLanguages {
			if targetLanguage.Language == input.TargetLanguage {
				targetLangApiName = targetLanguage.APIName
			}
		}
		if sourceLangApiName != "" && targetLangApiName != "" {
			url := dictionary.URL
			url = strings.Replace(url, "##SOURCE_LANGUAGE##", sourceLangApiName, 1)
			url = strings.Replace(url, "##TARGET_LANGUAGE##", targetLangApiName, 1)
			url = strings.Replace(url, "##LEMMA##", input.Lemma, 1)
			return url, nil
		}
		sourceLangApiName = ""
		targetLangApiName = ""
	}
	return "", fmt.Errorf("no dictionary found for source language %s and target language %s", input.SourceLanguage, input.TargetLanguage)
}

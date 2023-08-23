package types

import (
	"fmt"
	"strings"
)

type Language string

const (
	Bulgarian       Language = "bg"
	Czech           Language = "cs"
	Danish          Language = "da"
	German          Language = "de"
	Greek           Language = "el"
	English         Language = "en"
	Spanish         Language = "es"
	Estonian        Language = "et"
	Finnish         Language = "fi"
	French          Language = "fr"
	Hungarian       Language = "hu"
	Indonesian      Language = "id"
	Italian         Language = "it"
	Japanese        Language = "ja"
	Korean          Language = "ko"
	Lithuanian      Language = "lt"
	Latvian         Language = "lv"
	NorwegianBokmål Language = "nb"
	Dutch           Language = "nl"
	Polish          Language = "pl"
	Portuguese      Language = "pt"
	Romanian        Language = "ro"
	Russian         Language = "ru"
	Slovak          Language = "sk"
	Slovenian       Language = "sl"
	Swedish         Language = "sv"
	Turkish         Language = "tr"
	Ukrainian       Language = "uk"
	Chinese         Language = "zh"
)

var SupportedLanguages []Language = []Language{
	Bulgarian,
	Czech,
	Danish,
	German,
	Greek,
	English,
	Spanish,
	Estonian,
	Finnish,
	French,
	Hungarian,
	Indonesian,
	Italian,
	Japanese,
	Korean,
	Lithuanian,
	Latvian,
	NorwegianBokmål,
	Dutch,
	Polish,
	Portuguese,
	Romanian,
	Russian,
	Slovak,
	Slovenian,
	Swedish,
	Turkish,
	Ukrainian,
	Chinese,
}

func NewLanguage(inputLang string) (Language, error) {
	lang := Language(inputLang)
	for _, supportedLang := range SupportedLanguages {
		if lang == supportedLang {
			return lang, nil
		}
	}
	return "", fmt.Errorf("Language %s is not supported", inputLang)
}

func (lang Language) Uppercase() string {
	return strings.ToUpper(string(lang))
}

func (lang Language) Lowercase() string {
	return strings.ToLower(string(lang))
}

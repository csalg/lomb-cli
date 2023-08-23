package deepl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process/translators"
	"github.com/csalg/lomb-cli/pkg/types"
)

var DeeplSupportedLanguages []types.Language = []types.Language{
	types.Bulgarian,
	types.Czech,
	types.Danish,
	types.German,
	types.Greek,
	types.English,
	types.Spanish,
	types.Estonian,
	types.Finnish,
	types.French,
	types.Hungarian,
	types.Indonesian,
	types.Italian,
	types.Japanese,
	types.Korean,
	types.Lithuanian,
	types.Latvian,
	types.NorwegianBokmål,
	types.Dutch,
	types.Polish,
	types.Portuguese,
	types.Romanian,
	types.Russian,
	types.Slovak,
	types.Slovenian,
	types.Swedish,
	types.Turkish,
	types.Ukrainian,
	types.Chinese,
}

type DeeplTranslator struct {
	httpClient *http.Client
	apiKey     string
	apiPro     bool
}

func New(key string, pro bool) DeeplTranslator {
	return DeeplTranslator{
		httpClient: &http.Client{},
		apiKey:     key,
		apiPro:     pro,
	}
}

func IsLanguageSupported(lang types.Language) bool {
	for _, supportedLang := range DeeplSupportedLanguages {
		if lang == supportedLang {
			return true
		}
	}
	return false
}

func (dt DeeplTranslator) Translate(sourceLang, targetLang types.Language, text []string) ([]translators.TranslatedText, error) {
	url := "https://api-free.deepl.com/v2/translate"
	if dt.apiPro {
		url = "https://api.deepl.com/v2/translate"
	}
	method := "POST"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return []translators.TranslatedText{}, fmt.Errorf("creating request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", dt.apiKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("source_lang", sourceLang.Uppercase())
	q.Add("target_lang", targetLang.Uppercase())
	for _, t := range text {
		q.Add("text", t)
	}
	req.URL.RawQuery = q.Encode()

	res, err := dt.httpClient.Do(req)
	if err != nil {
		return []translators.TranslatedText{}, fmt.Errorf("doing request: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []translators.TranslatedText{}, fmt.Errorf("reading response: %s", err)
	}

	var parsedRes DeeplResponse
	err = json.Unmarshal(body, &parsedRes)
	if err != nil {
		return []translators.TranslatedText{}, fmt.Errorf("unmarshalling response: %s", err)
	}

	if len(parsedRes.Translations) != len(text) {
		return []translators.TranslatedText{}, fmt.Errorf("unexpected number of translations: expected %d, got %d", len(text), len(parsedRes.Translations))
	}

	translatedTexts := []translators.TranslatedText{}
	for i, t := range parsedRes.Translations {
		translatedTexts = append(translatedTexts, translators.TranslatedText{
			Translated: t.Text,
			Original:   text[i],
		})
	}
	return translatedTexts, nil
}

type DeeplResponse struct {
	Translations []DeeplResponseItem `json:"translations"`
}

type DeeplResponseItem struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

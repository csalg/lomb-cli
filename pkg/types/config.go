package types

import "fmt"

type Config struct {
	AppConfig
	*LanguagePairConfig
}

type AppConfig struct {
	OpenAIAPIKey          string `json:"openai_api_key,omitempty"`
	DeeplAPIKey           string `json:"deepl_api_key,omitempty"`
	DeeplAPIPro           bool   `json:"deepl_api_pro,omitempty"`
	GoogleTranslateAPIKey string `json:"google_translate_api_key,omitempty"`
	Port                  string `json:"port,omitempty"`
	BaseURL               string `json:"base_url,omitempty"`
}

// ClientURL returns the client URL (e.g. http://localhost:9090)
func (c AppConfig) ClientURL() string {
	return fmt.Sprintf("%s:%s", c.BaseURL, c.Port)
}

type LanguagePairConfig struct {
	TargetLanguage string
	BaseLanguage   string
	Views          []View
}

type View struct {
	Name string
	Grid [][]Cell
}

type Cell struct {
	Type          CellType
	ID            string
	DictionaryURL string
}

type CellType string

const (
	CellLemmaCounts             CellType = "lemma-counts"
	CellExamples                CellType = "examples"
	CellDictionary              CellType = "dictionary"
	CellUnderstandableSentences CellType = "understandable-sentences"
	CellReader                  CellType = "reader"
)

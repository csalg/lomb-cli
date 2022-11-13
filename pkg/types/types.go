package types

// Config is the configuration for the lomb CLI.
type Config struct {
	Dictionaries []Dictionary `json:"dictionaries"`
	Port         int          `json:"port"`
	Host         string       `json:"host"`
}

type Dictionary struct {
	Name            string                     `json:"name"`
	URL             string                     `json:"url"`
	SourceLanguages []LanguageDictionaryConfig `json:"source_languages"`
	TargetLanguages []LanguageDictionaryConfig `json:"target_languages"`
}

type LanguageDictionaryConfig struct {
	Language string `json:"language"`
	APIName  string `json:"api_name"`
}

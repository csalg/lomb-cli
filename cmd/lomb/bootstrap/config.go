package bootstrap

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	OpenAIAPIKey          string `json:"openai_api_key,omitempty"`
	DeeplAPIKey           string `json:"deepl_api_key,omitempty"`
	DeeplAPIPro           bool   `json:"deepl_api_pro,omitempty"`
	GoogleTranslateAPIKey string `json:"google_translate_api_key,omitempty"`
	Port                  string `json:"port,omitempty"`
	BaseURL               string `json:"base_url,omitempty"`
}

func ReadConfig() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("getting user config directory: %w", err)
	}
	filename := filepath.Join(configDir, "lomb", "config.json")
	if _, err := os.Stat(filename); err != nil {
		return Config{}, fmt.Errorf("statting config file: %w", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return Config{}, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return Config{}, fmt.Errorf("unmarshaling config file: %w", err)
	}
	if conf.Port == "" {
		conf.Port = "9090"
	}
	if conf.BaseURL == "" {
		conf.BaseURL = "http://localhost"
	}

	return conf, nil
}

// ClientURL returns the client URL (e.g. http://localhost:9090)
func (c Config) ClientURL() string {
	return fmt.Sprintf("%s:%s", c.BaseURL, c.Port)
}
